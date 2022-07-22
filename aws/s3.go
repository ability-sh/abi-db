package aws

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"

	"github.com/ability-sh/abi-db/source"
	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func init() {
	source.Reg("aws-s3", NewS3)
}

type s3Source struct {
	client *s3.Client
	bucket string
}

func Unwrap(err error) error {
	for {
		e := errors.Unwrap(err)
		if e == nil {
			return err
		}
		err = e
	}
	return err
}

func NewS3(driver string, sConfig interface{}) (source.Source, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(dynamic.StringValue(dynamic.Get(sConfig, "region"), "ap-northeast-1")),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     dynamic.StringValue(dynamic.Get(sConfig, "accesskey"), ""),
				SecretAccessKey: dynamic.StringValue(dynamic.Get(sConfig, "secret"), ""),
			},
		}),
	)

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &s3Source{client: client, bucket: dynamic.StringValue(dynamic.Get(sConfig, "bucket"), "")}, nil
}

func (s *s3Source) Put(key string, data []byte) error {
	ctx := context.Background()
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{Bucket: &s.bucket, Key: &key, Body: bytes.NewReader(data)})
	return err
}

func (s *s3Source) Del(key string) error {
	ctx := context.Background()
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: &s.bucket, Key: &key})
	return err
}

func (s *s3Source) Get(key string) ([]byte, error) {
	ctx := context.Background()
	rs, err := s.client.GetObject(ctx, &s3.GetObjectInput{Bucket: &s.bucket, Key: &key})
	if err != nil {
		_, ok := Unwrap(err).(*types.NoSuchKey)
		if ok {
			return nil, source.ErrNoSuchKey
		}
		return nil, err
	}
	defer rs.Body.Close()
	return ioutil.ReadAll(rs.Body)
}

type s3Cursor struct {
	s           *s3Source
	prefix      string
	delimiter   string
	isTruncated bool
	contents    []types.Object
	index       int
	ctx         context.Context
}

func (s *s3Cursor) Next() (string, error) {
	if s.contents == nil {
		rs, err := s.s.client.ListObjectsV2(s.ctx, &s3.ListObjectsV2Input{Bucket: &s.s.bucket, Prefix: &s.prefix})
		if err != nil {
			return "", err
		}
		s.delimiter = *rs.Delimiter
		s.isTruncated = rs.IsTruncated
		s.contents = rs.Contents
	}
	if s.index >= len(s.contents) {
		if s.isTruncated {
			rs, err := s.s.client.ListObjectsV2(s.ctx, &s3.ListObjectsV2Input{Bucket: &s.s.bucket, Prefix: &s.prefix, Delimiter: &s.delimiter})
			if err != nil {
				return "", err
			}
			s.delimiter = *rs.Delimiter
			s.isTruncated = rs.IsTruncated
			s.contents = rs.Contents
		} else {
			return "", io.EOF
		}
	}
	if s.index < len(s.contents) {
		r := s.contents[s.index]
		s.index += 1
		return *r.Key, nil
	}
	return "", io.EOF
}

func (s *s3Cursor) Close() {

}

func (s *s3Source) Query(prefix string) (source.Cursor, error) {
	return &s3Cursor{s: s, prefix: prefix, ctx: context.Background()}, nil

}
