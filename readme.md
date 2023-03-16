# Get storage classes and size of all AWS S3 buckets/objects in a region

If the bucket contains more than 100k object, you will see a message counting up.
Pages will be counted in parallel.

With that data you perform a pricing analysis for s3 storage class pricing.

Example: 

```bash
go run sac
Bucket: amplify-awsamplifyauthstarte-dev-155522-deployment
STANDARD 8
44 kB
---
Bucket: amplify-totpcognito-dev-112355-deployment
0 B
---
Bucket: amplify-totpcognito-dev-112613-deployment
0 B
---
Bucket: amplify-totplogin-dev-160212-deployment
STANDARD 7
14 kB
---
Bucket: amplify-trainerportal-dev-90853-deployment
Bucket not in region eu-central-1, but in region eu-west-1
---
Bucket: aws-cloudtrail-logs-123456789012-d2d863a5
100 k,  200 k,  300 k,  400 k,  500 k,  600 k,  700 k,  800 k, etc
```
