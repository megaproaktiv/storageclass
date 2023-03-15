package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/dustin/go-humanize"
)

var client *s3.Client
var uploader *manager.Uploader
var region string

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client = s3.NewFromConfig(cfg)
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

func storeageclass(client *s3.Client, bucketName *string) (map[string]int, uint64) {
	var size uint64
	c := Container{
        counters: map[string]int{},
    }
	maxKeys := 1000
	params := &s3.ListObjectsV2Input{
		Bucket: bucketName,
	}
	// Create the Paginator for the ListObjectsV2 operation.
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
		if (modulo == 0){
			fmt.Printf(" %v k, ", i)
		}
		
		
		// Next Page takes a new context for each page retrieval. This is where
		// you could add timeouts or deadlines.
		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
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

func main() {
	resp, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		fmt.Println("Cannot list buckets")
		panic(err)
	}
	for _, buck := range resp.Buckets {
		name := buck.Name
		fmt.Printf("Bucket: %v\n", *name)
		
		// find region
		respRegion, err := client.GetBucketLocation(context.TODO(), &s3.GetBucketLocationInput{
			Bucket:              name,
		})
		if err != nil {
			fmt.Println("Cannot get bucket region ")
			panic(err)
		}
		if respRegion.LocationConstraint == types.BucketLocationConstraint(region){	
			classes,size := storeageclass(client, name)
			for key, value := range classes {
				fmt.Println(key, value)
			}
			fmt.Printf("%s\n", humanize.Bytes(size))

		}else{
			fmt.Printf("Bucket not in region %v, but in region %v\n",region, respRegion.LocationConstraint)
		}
		fmt.Println("---")
	}
}
