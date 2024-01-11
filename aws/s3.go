package aws

import (
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/skit-ai/vcore/errors"
	"github.com/skit-ai/vcore/log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	// Regex for S3 URLs, VPCE interface endpoint
	vpceURLPattern = "^((.+)\\.)?" + // maybe a bucket name
		"(bucket|accesspoint|control)\\.vpce-[-a-z0-9]+\\." + // VPC endpoint DNS name
		"s3[.-]" + // S3 service name
		"(([-a-z0-9]+)\\.)?" + // region name, optional for us-east-1
		"vpce\\." +
		"(amazonaws\\.com|c2s\\.ic\\.gov|sc2s\\.sgov\\.gov)"
	vpceURLPatternHostIdx   = 0
	vpceURLPatternBucketIdx = 2
	vpceURLPatternRegionIdx = 5

	// Regex for S3 URLs, public S3 endpoint
	nonVpceURLPattern = "^((.+)\\.)?" + // maybe a bucket name
		"s3[.-](website[-.])?(accelerate\\.)?(dualstack[-.])?" + // S3 service name with optional features
		"(([-a-z0-9]+)\\.)?" + // region name, optional for us-east-1
		"(amazonaws\\.com|c2s\\.ic\\.gov|sc2s\\.sgov\\.gov)"
	nonVpceURLPatternBucketIdx = 2
	nonVpceURLPatternRegionIdx = 7
)

var (
	vpceUrlRegex    = regexp.MustCompile(vpceURLPattern)
	nonVpceUrlRegex = regexp.MustCompile(nonVpceURLPattern)
)

// S3URL holds interesting pieces after parsing a s3 URL
type S3URL struct {
	IsPathStyle bool
	EndPoint    string
	Bucket      string
	Key         string
	Region      string
}

// DownloadFile downloads a file from s3 based on the key and writes it into WriteAt.
func (u S3URL) DownloadFile(ctx context.Context, w io.WriterAt) error {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(u.Region), // Specify the region where the bucket is located
		Endpoint: aws.String(u.EndPoint),
	})
	if err != nil {
		return errors.NewError("Error creating session", err, false)
	}

	downloader := s3manager.NewDownloader(sess)

	numBytes, err := downloader.DownloadWithContext(ctx, w, &s3.GetObjectInput{
		Bucket: aws.String(u.Bucket),
		Key:    aws.String(u.Key),
	})

	if err != nil {
		return errors.NewError("Error downloading file", err, false)
	}

	slog.Debug("Downloaded file", "size", numBytes)

	return nil
}

// ParseAmazonS3URL parses an HTTP/HTTPS URL for an S3 resource and returns an
// S3URL object.
//
// S3 URLs come in two flavors: virtual hosted-style URLs and path-style URLs.
// Virtual hosted-style URLs have the bucket name as the first component of the
// hostname, e.g.
//
//	https://mybucket.s3.us-east-1.amazonaws.com/a/b/c
//
// Path-style URLs have the bucket name as the first component of the path, e.g.
//
//	https://s3.us-east-1.amazonaws.com/mybucket/a/b/c
func ParseAmazonS3URL(s3URL *url.URL) (S3URL, error) {
	output, err := parseBucketAndRegionFromHost(s3URL.Host)
	if err != nil {
		return S3URL{}, errors.NewError("parsing host failed", err, false)
	}

	output.IsPathStyle = output.Bucket == ""

	path := s3URL.Path

	if output.IsPathStyle {
		// no bucket name in the authority, parse it from the path
		output.IsPathStyle = true

		// grab the encoded path so we don't run afoul of '/'s in the bucket name
		if path == "/" || path == "" {
		} else {
			path = path[1:]
			index := strings.Index(path, "/")
			if index == -1 {
				// https://s3.amazonaws.com/bucket
				output.Bucket = path
				output.Key = ""
			} else if index == (len(path) - 1) {
				// https://s3.amazonaws.com/bucket/
				output.Bucket = strings.TrimRight(path, "/")
				output.Key = ""
			} else {
				// https://s3.amazonaws.com/bucket/key
				output.Bucket = path[:index]
				output.Key = path[index+1:]
			}
		}
	} else {
		// bucket name in the host, path is the object key
		if path == "/" || path == "" {
			output.Key = ""
		} else {
			output.Key = path[1:]
		}
	}

	if strings.EqualFold(output.Region, "external-1") {
		output.Region = "us-east-1"
	} else if output.Region == "" {
		// s3 bucket URL in us-east-1 doesn't include region
		output.Region = "us-east-1"
	}

	return output, nil
}

func parseBucketAndRegionFromHost(host string) (S3URL, error) {
	result := vpceUrlRegex.FindStringSubmatch(host)
	if result != nil && len(result) > vpceURLPatternBucketIdx && len(result) > vpceURLPatternRegionIdx {
		return S3URL{
			EndPoint: result[vpceURLPatternHostIdx],
			Bucket:   result[vpceURLPatternBucketIdx],
			Region:   result[vpceURLPatternRegionIdx],
		}, nil
	} else {
		result = nonVpceUrlRegex.FindStringSubmatch(host)
		if result != nil && len(result) > vpceURLPatternBucketIdx && len(result) > vpceURLPatternRegionIdx {
			return S3URL{
				Bucket: result[nonVpceURLPatternBucketIdx],
				Region: result[nonVpceURLPatternRegionIdx],
			}, nil
		} else {
			return S3URL{}, errors.NewError("failed to match URL", nil, false)
		}
	}
}

// DownloadFileFromS3 takes an S3 URL and a filePath, downloads the file from s3 and stores it in the filePath.
func DownloadFileFromS3(ctx context.Context, downloadURL, filePath string) error {
	parsedURL, err := url.Parse(downloadURL)
	if err != nil {
		return err
	}

	// Parse s3 URL to extract region, key and bucket.
	s3URL, err := ParseAmazonS3URL(parsedURL)
	if err != nil {
		return errors.NewError("Failed to parse URL", err, false)
	}

	// Create file path
	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return errors.NewError("Unable to create directory", err, false)
	}

	// Create a local file to write to
	f, err := os.Create(filePath)
	if err != nil {
		return errors.NewError("Error creating file", err, false)
	}

	defer func() {
		// Ensure file is closed even if an error occurs
		if f != nil {
			f.Close()
		}
	}()

	return s3URL.DownloadFile(ctx, f)
}
