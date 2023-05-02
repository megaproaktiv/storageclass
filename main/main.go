package main

import (
	"context"
	"fmt"
	"sac"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dustin/go-humanize"
)



func main() {
	resp, err := sac.Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		fmt.Println("Cannot list buckets")
		panic(err)
	}
	var sumSize uint64
	for _, buck := range resp.Buckets {
		name := buck.Name
		
		// find region
		respRegion, err := sac.Client.GetBucketLocation(context.TODO(), &s3.GetBucketLocationInput{
			Bucket: name,
		})
		if err != nil {
			fmt.Println("Cannot get bucket region ")
			panic(err)
		}
		region := respRegion.LocationConstraint
		fmt.Printf("Region/Bucket:%v -  %v\n", region,*name)

		classes, size := sac.Storeageclass( name, string(region))
		for key, value := range classes {
			fmt.Println(key, value)
		}
		fmt.Printf("%s\n", humanize.Bytes(size))
		sumSize += size

		fmt.Println("---")
	}
	fmt.Printf("Size of all Buckets: %v\n", humanize.Bytes(sumSize))
}
