Simple downloader from AWS Glacier
==================================

Could not find simple utility to download archive from AWS Glacier to my PC,
so written my own.

Usage
-----

```
aws configure
aws glacier describe-job --account-id - --vault-name YOUR_VOULT_NAME --job-id YOUR_LONG_JOB_ID > job.json
./gdown job.json OUTPUT_FILE
```

Where `YOUR_VOULT_NAME`, `YOUR_LONG_JOB_ID` and `OUTPUT_FILE` are variables

Assumptions
-----------

* You have AWS CLI configured with Glacer permisions
* You have something uploaded to Glacer `aws glacier upload-archive ...`
* You have compleated retrieval Job request `aws glacier initiate-job ...`

Known issues
------------

* This is fast-written script, that just does its work. It was not intended to be used a library.
* Will not work on AWS EC2: `Error: open /home/ec2-user/.aws/credentials: no such file or directory`

Compiling
---------

```
go get github.com/smartystreets/go-aws-auth
go build -o gdown main.go
```

Useful for development
```
go run main.go job.json some-file.gzip
```

References/credits
------------------

* [SmartyStreets for go-aws-auth](https://github.com/smartystreets/go-aws-auth)
* [AWS: How to initiate a job](http://docs.aws.amazon.com/cli/latest/reference/glacier/initiate-job.html)
* [AWS: Descriptions for Tiers/price](http://docs.aws.amazon.com/amazonglacier/latest/dev/api-initiate-job-post.html)

Licence
-------

MIT