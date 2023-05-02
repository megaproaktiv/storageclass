package sac

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var Client *s3.Client
var uploader *manager.Uploader
var region string

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	Client = s3.NewFromConfig(cfg)
	region = cfg.Region

}

type Container struct {
	mu       sync.Mutex
	counters map[string]int
}

func (c *Container) inc(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counters[name]++
}

func Storeageclass(bucketName *string, region string) (map[string]int, uint64) {
	var size uint64
	c := Container{
		counters: map[string]int{},
	}
	maxKeys := 1000
	params := &s3.ListObjectsV2Input{
		Bucket: bucketName,
	}
	if region == "" {
		region = "us-east-1"
	}
	// Create the Paginator for the ListObjectsV2 operation
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := s3.NewFromConfig(cfg)
	p := s3.NewListObjectsV2Paginator(client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
		if v := int32(maxKeys); v != 0 {
			o.Limit = v
		}
	})

	// Iterate through the S3 object pages, printing each object returned.
	var i int
	notify := 100
	var wg sync.WaitGroup

	for p.HasMorePages() {
		i++
		modulo := i % notify
		if modulo == 0 {
			fmt.Printf(" %v k, ", i)
		}

		// Next Page takes a new context for each page retrieval. This is where
		// you could add timeouts or deadlines.
		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Printf("failed to get page %v, %v", i, err)
			return c.counters, size
		}

		wg.Add(1)

		go func(objects []types.Object) {
			defer wg.Done()

			for _, obj := range page.Contents {
				// fmt.Println("Class:", obj.StorageClass)
				c.inc(string(obj.StorageClass))
				size += uint64(obj.Size)
			}

		}(page.Contents)

		// Log the objects found
		wg.Wait()
	}
	return c.counters, size
}
