package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var minioCmd = &cobra.Command{
	Use:   "minio",
	Short: "MinIO object storage management commands",
	Long:  `Manage MinIO buckets, objects, and policies`,
}

var minioListBucketsCmd = &cobra.Command{
	Use:   "buckets",
	Short: "List all buckets",
	Run: func(cmd *cobra.Command, args []string) {
		listBuckets()
	},
}

var minioCreateBucketCmd = &cobra.Command{
	Use:   "create-bucket [bucket-name]",
	Short: "Create a new bucket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		createBucket(args[0], region)
	},
}

var minioDeleteBucketCmd = &cobra.Command{
	Use:   "delete-bucket [bucket-name]",
	Short: "Delete a bucket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		deleteBucket(args[0], force)
	},
}

var minioListObjectsCmd = &cobra.Command{
	Use:   "list [bucket-name]",
	Short: "List objects in a bucket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prefix, _ := cmd.Flags().GetString("prefix")
		recursive, _ := cmd.Flags().GetBool("recursive")
		listObjects(args[0], prefix, recursive)
	},
}

var minioUploadCmd = &cobra.Command{
	Use:   "upload [bucket-name] [local-file] [object-name]",
	Short: "Upload a file to a bucket",
	Args:  cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		objectName := ""
		if len(args) == 3 {
			objectName = args[2]
		} else {
			objectName = filepath.Base(args[1])
		}
		uploadFile(args[0], args[1], objectName)
	},
}

var minioDownloadCmd = &cobra.Command{
	Use:   "download [bucket-name] [object-name] [local-file]",
	Short: "Download an object from a bucket",
	Args:  cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		localFile := ""
		if len(args) == 3 {
			localFile = args[2]
		} else {
			localFile = filepath.Base(args[1])
		}
		downloadFile(args[0], args[1], localFile)
	},
}

var minioDeleteObjectCmd = &cobra.Command{
	Use:   "delete [bucket-name] [object-name]",
	Short: "Delete an object from a bucket",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		deleteObject(args[0], args[1])
	},
}

var minioCopyCmd = &cobra.Command{
	Use:   "copy [source-bucket] [source-object] [dest-bucket] [dest-object]",
	Short: "Copy an object between buckets",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		copyObject(args[0], args[1], args[2], args[3])
	},
}

var minioStatCmd = &cobra.Command{
	Use:   "stat [bucket-name] [object-name]",
	Short: "Get object information",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			statBucket(args[0])
		} else {
			statObject(args[0], args[1])
		}
	},
}

var minioMirrorCmd = &cobra.Command{
	Use:   "mirror [local-dir] [bucket-name]",
	Short: "Mirror local directory to bucket",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		prefix, _ := cmd.Flags().GetString("prefix")
		mirrorDirectory(args[0], args[1], prefix)
	},
}

func init() {
	minioCmd.PersistentFlags().StringP("endpoint", "e", "", "MinIO endpoint (env: MINIO_ENDPOINT)")
	minioCmd.PersistentFlags().StringP("access-key", "a", "", "Access key (env: MINIO_ACCESS_KEY)")
	minioCmd.PersistentFlags().StringP("secret-key", "s", "", "Secret key (env: MINIO_SECRET_KEY)")
	minioCmd.PersistentFlags().BoolP("use-ssl", "S", false, "Use SSL (env: MINIO_USE_SSL)")

	minioCreateBucketCmd.Flags().StringP("region", "r", "us-east-1", "Bucket region")
	minioDeleteBucketCmd.Flags().BoolP("force", "f", false, "Force delete non-empty bucket")
	minioListObjectsCmd.Flags().StringP("prefix", "p", "", "Object prefix")
	minioListObjectsCmd.Flags().BoolP("recursive", "r", false, "List recursively")
	minioMirrorCmd.Flags().StringP("prefix", "p", "", "Bucket prefix")

	minioCmd.AddCommand(minioListBucketsCmd)
	minioCmd.AddCommand(minioCreateBucketCmd)
	minioCmd.AddCommand(minioDeleteBucketCmd)
	minioCmd.AddCommand(minioListObjectsCmd)
	minioCmd.AddCommand(minioUploadCmd)
	minioCmd.AddCommand(minioDownloadCmd)
	minioCmd.AddCommand(minioDeleteObjectCmd)
	minioCmd.AddCommand(minioCopyCmd)
	minioCmd.AddCommand(minioStatCmd)
	minioCmd.AddCommand(minioMirrorCmd)
}

func getMinIOClient() (*minio.Client, error) {
	endpoint, _ := minioCmd.Flags().GetString("endpoint")
	accessKey, _ := minioCmd.Flags().GetString("access-key")
	secretKey, _ := minioCmd.Flags().GetString("secret-key")
	useSSL, _ := minioCmd.Flags().GetBool("use-ssl")

	// Use environment variables if flags not provided
	if endpoint == "" {
		endpoint = viper.GetString("minio.endpoint")
	}
	if accessKey == "" {
		accessKey = viper.GetString("minio.access_key")
	}
	if secretKey == "" {
		secretKey = viper.GetString("minio.secret_key")
	}
	if !useSSL {
		useSSL = viper.GetBool("minio.use_ssl")
	}

	// Initialize minio client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

func listBuckets() {
	client, err := getMinIOClient()
	if err != nil {
		color.Red("Error connecting to MinIO: %v", err)
		return
	}

	ctx := context.Background()
	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		color.Red("Error listing buckets: %v", err)
		return
	}

	color.Green("Buckets:")
	for _, bucket := range buckets {
		fmt.Printf("  - %s (created: %s)\n", bucket.Name, bucket.CreationDate.Format("2006-01-02 15:04:05"))
	}
}

func createBucket(bucketName, region string) {
	client, err := getMinIOClient()
	if err != nil {
		color.Red("Error connecting to MinIO: %v", err)
		return
	}

	ctx := context.Background()
	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: region})
	if err != nil {
		// Check if bucket already exists
		exists, errBucketExists := client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			color.Yellow("Bucket '%s' already exists", bucketName)
		} else {
			color.Red("Error creating bucket: %v", err)
		}
		return
	}

	color.Green("Bucket '%s' created successfully", bucketName)
}

func deleteBucket(bucketName string, force bool) {
	client, err := getMinIOClient()
	if err != nil {
		color.Red("Error connecting to MinIO: %v", err)
		return
	}

	ctx := context.Background()

	if force {
		// Remove all objects in bucket first
		objectsCh := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
			Recursive: true,
		})

		for object := range objectsCh {
			if object.Err != nil {
				color.Red("Error listing objects: %v", object.Err)
				return
			}

			err := client.RemoveObject(ctx, bucketName, object.Key, minio.RemoveObjectOptions{})
			if err != nil {
				color.Red("Error removing object '%s': %v", object.Key, err)
				return
			}
		}
	}

	err = client.RemoveBucket(ctx, bucketName)
	if err != nil {
		color.Red("Error deleting bucket: %v", err)
		return
	}

	color.Green("Bucket '%s' deleted successfully", bucketName)
}

func listObjects(bucketName, prefix string, recursive bool) {
	client, err := getMinIOClient()
	if err != nil {
		color.Red("Error connecting to MinIO: %v", err)
		return
	}

	ctx := context.Background()

	objectsCh := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	})

	color.Green("Objects in bucket '%s':", bucketName)
	count := 0
	totalSize := int64(0)
	for object := range objectsCh {
		if object.Err != nil {
			color.Red("Error: %v", object.Err)
			return
		}
		fmt.Printf("  - %s (size: %d bytes, modified: %s)\n",
			object.Key, object.Size, object.LastModified.Format("2006-01-02 15:04:05"))
		count++
		totalSize += object.Size
	}
	fmt.Printf("\nTotal: %d objects, %d bytes\n", count, totalSize)
}

func uploadFile(bucketName, localFile, objectName string) {
	client, err := getMinIOClient()
	if err != nil {
		color.Red("Error connecting to MinIO: %v", err)
		return
	}

	ctx := context.Background()

	// Check if file exists
	file, err := os.Open(localFile)
	if err != nil {
		color.Red("Error opening file: %v", err)
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		color.Red("Error getting file info: %v", err)
		return
	}

	// Upload the file
	_, err = client.PutObject(ctx, bucketName, objectName, file, fileStat.Size(), minio.PutObjectOptions{})
	if err != nil {
		color.Red("Error uploading file: %v", err)
		return
	}

	color.Green("File '%s' uploaded successfully as '%s'", localFile, objectName)
}

func downloadFile(bucketName, objectName, localFile string) {
	client, err := getMinIOClient()
	if err != nil {
		color.Red("Error connecting to MinIO: %v", err)
		return
	}

	ctx := context.Background()

	// Get object from bucket
	object, err := client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		color.Red("Error getting object: %v", err)
		return
	}
	defer object.Close()

	// Create local file
	file, err := os.Create(localFile)
	if err != nil {
		color.Red("Error creating local file: %v", err)
		return
	}
	defer file.Close()

	// Copy object to file
	_, err = io.Copy(file, object)
	if err != nil {
		color.Red("Error downloading file: %v", err)
		return
	}

	color.Green("Object '%s' downloaded successfully as '%s'", objectName, localFile)
}

func deleteObject(bucketName, objectName string) {
	client, err := getMinIOClient()
	if err != nil {
		color.Red("Error connecting to MinIO: %v", err)
		return
	}

	ctx := context.Background()

	err = client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		color.Red("Error deleting object: %v", err)
		return
	}

	color.Green("Object '%s' deleted successfully from bucket '%s'", objectName, bucketName)
}

func copyObject(sourceBucket, sourceObject, destBucket, destObject string) {
	client, err := getMinIOClient()
	if err != nil {
		color.Red("Error connecting to MinIO: %v", err)
		return
	}

	ctx := context.Background()

	// Copy object
	srcOpts := minio.CopySrcOptions{
		Bucket: sourceBucket,
		Object: sourceObject,
	}

	dstOpts := minio.CopyDestOptions{
		Bucket: destBucket,
		Object: destObject,
	}

	_, err = client.CopyObject(ctx, dstOpts, srcOpts)
	if err != nil {
		color.Red("Error copying object: %v", err)
		return
	}

	color.Green("Object copied successfully from '%s/%s' to '%s/%s'",
		sourceBucket, sourceObject, destBucket, destObject)
}

func statBucket(bucketName string) {
	client, err := getMinIOClient()
	if err != nil {
		color.Red("Error connecting to MinIO: %v", err)
		return
	}

	ctx := context.Background()

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		color.Red("Error checking bucket: %v", err)
		return
	}

	if !exists {
		color.Yellow("Bucket '%s' does not exist", bucketName)
		return
	}

	// Count objects and total size
	objectsCh := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	count := 0
	totalSize := int64(0)
	for object := range objectsCh {
		if object.Err != nil {
			color.Red("Error: %v", object.Err)
			return
		}
		count++
		totalSize += object.Size
	}

	color.Green("Bucket '%s' information:", bucketName)
	fmt.Printf("  - Objects: %d\n", count)
	fmt.Printf("  - Total size: %d bytes\n", totalSize)
}

func statObject(bucketName, objectName string) {
	client, err := getMinIOClient()
	if err != nil {
		color.Red("Error connecting to MinIO: %v", err)
		return
	}

	ctx := context.Background()

	info, err := client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		color.Red("Error getting object info: %v", err)
		return
	}

	color.Green("Object '%s' information:", objectName)
	fmt.Printf("  - Bucket: %s\n", bucketName)
	fmt.Printf("  - Size: %d bytes\n", info.Size)
	fmt.Printf("  - Last modified: %s\n", info.LastModified.Format("2006-01-02 15:04:05"))
	fmt.Printf("  - ETag: %s\n", info.ETag)
	fmt.Printf("  - Content-Type: %s\n", info.ContentType)
}

func mirrorDirectory(localDir, bucketName, prefix string) {
	client, err := getMinIOClient()
	if err != nil {
		color.Red("Error connecting to MinIO: %v", err)
		return
	}

	ctx := context.Background()

	// Walk through local directory
	uploadCount := 0
	err = filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return err
		}

		// Convert to forward slashes for object name
		objectName := strings.ReplaceAll(relPath, "\\", "/")
		if prefix != "" {
			objectName = prefix + "/" + objectName
		}

		// Upload file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = client.PutObject(ctx, bucketName, objectName, file, info.Size(), minio.PutObjectOptions{})
		if err != nil {
			color.Red("Error uploading '%s': %v", path, err)
			return nil
		}

		fmt.Printf("Uploaded: %s -> %s/%s\n", path, bucketName, objectName)
		uploadCount++
		return nil
	})

	if err != nil {
		color.Red("Error walking directory: %v", err)
		return
	}

	color.Green("Mirror completed: %d files uploaded", uploadCount)
}

func GetMinIOCommand() *cobra.Command {
	return minioCmd
}
