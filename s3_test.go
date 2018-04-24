package s3

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHelper(t *testing.T) {
	config := Config{
		AccessKeyID:     "",
		Endpoint:        "",
		Region:          "",
		SecretAccessKey: "",
		BucketName:      "",
		SSL:             false,
	}
	Convey("New", t, func() {
		Convey("Config validation error", func() {
			s3, err := New(config)
			// config validation error
			So(err, ShouldNotBeNil)
			So(s3, ShouldBeNil)
		})

		Convey("NewWithRegion error", func() {
			// minio.NewWithRegion throws error if endpoint is an invalid host
			config.Endpoint = "invalid:host:xxx"
			config.AccessKeyID = "x"
			config.SecretAccessKey = "x"
			config.Region = "x"
			config.BucketName = "x"
			s3, err := New(config)

			So(err, ShouldNotBeNil)
			So(s3, ShouldBeNil)
		})

		Convey("Success", func() {
			config.Endpoint = "localhost"
			config.AccessKeyID = "x"
			config.SecretAccessKey = "x"
			config.Region = "x"
			config.BucketName = "x"
			s3, err := New(config)

			So(err, ShouldBeNil)
			So(s3, ShouldNotBeNil)
		})
	})

	Convey("CreateBucket", t, func() {
		config := Config{
			AccessKeyID:     "x",
			Endpoint:        "localhost",
			Region:          "x",
			SecretAccessKey: "x",
			BucketName:      "x",
			SSL:             false,
		}

		Convey("Invalid bucket name", func() {
			s3, err := New(config)
			So(err, ShouldBeNil)
			// invalid bucket name (too short), minio makebucket throws error
			err = s3.CreateBucket("x")
			So(err, ShouldNotBeNil)
		})

		Convey("Directory created", func() {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Hello world!")
			}))

			url := strings.TrimPrefix(server.URL, "http://")
			config := Config{
				AccessKeyID:     "x",
				Endpoint:        url,
				Region:          "x",
				SecretAccessKey: "x",
				BucketName:      "x",
				SSL:             false,
			}

			s3, err := New(config)
			So(err, ShouldBeNil)
			err = s3.CreateBucket("x43563")
			So(err, ShouldBeNil)
		})

		Convey("Disabled S3", func() {
			s3 := helper{
				Enabled: false,
			}

			err := s3.CreateBucket("x")
			So(err, ShouldBeNil)
		})

	})

	Convey("CreateDirectory", t, func() {
		Convey("Disabled S3", func() {
			s3 := helper{
				Enabled: false,
			}

			err := s3.CreateDirectory("x", "asd")
			So(err, ShouldBeNil)
		})

		Convey("PutObject", func() {
			Convey("Success", func() {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, "{}")
				}))
				url := strings.TrimPrefix(server.URL, "http://")
				config := Config{
					AccessKeyID:     "x",
					Endpoint:        url,
					Region:          "x",
					SecretAccessKey: "x",
					BucketName:      "x",
					SSL:             false,
				}

				s3, err := New(config)
				So(err, ShouldBeNil)

				bucket := "x43563"
				err = s3.CreateDirectory(bucket, "1234678")
				So(err, ShouldBeNil)
			})
			Convey("Fail SetBucketPolicy", func() {
				i := 0
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if i > 0 {
						w.WriteHeader(400)
					}
					fmt.Fprintln(w, "{}")
					i++
				}))
				url := strings.TrimPrefix(server.URL, "http://")
				config := Config{
					AccessKeyID:     "x",
					Endpoint:        url,
					Region:          "x",
					SecretAccessKey: "x",
					BucketName:      "x",
					SSL:             false,
				}

				s3, err := New(config)
				So(err, ShouldBeNil)

				bucket := "x43563"
				err = s3.CreateDirectory(bucket, "1234678")
				So(err, ShouldNotBeNil)
			})

			Convey("Error", func() {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(400)
					fmt.Fprintln(w, "{\"error\": \"Invalid blabla\"}")
				}))

				url := strings.TrimPrefix(server.URL, "http://")
				config := Config{
					AccessKeyID:     "x",
					Endpoint:        url,
					Region:          "x",
					SecretAccessKey: "x",
					BucketName:      "x",
					SSL:             false,
				}

				s3, err := New(config)
				So(err, ShouldBeNil)
				bucket := "x43563"

				err = s3.CreateDirectory(bucket, "1234678")
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("CreateFile", t, func() {
		Convey("Disabled S3", func() {
			s3 := helper{
				Enabled: false,
			}
			bucket := "string"
			directory := "string"
			fileName := "string"
			content := bytes.NewReader([]byte("asdf"))
			length := int64(60)
			mime := "string"
			err := s3.CreateFile(bucket, directory, fileName, content, length, mime)
			So(err, ShouldBeNil)
		})
		Convey("Success", func() {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "{}")
			}))

			url := strings.TrimPrefix(server.URL, "http://")
			config := Config{
				AccessKeyID:     "x",
				Endpoint:        url,
				Region:          "x",
				SecretAccessKey: "x",
				BucketName:      "x",
				SSL:             false,
			}
			s3, err := New(config)
			So(err, ShouldBeNil)

			bucket := "string"
			directory := "string"
			fileName := "string.png"
			content := bytes.NewReader([]byte("asdf"))
			length := content.Len()
			mime := "image/png"
			err = s3.CreateFile(bucket, directory, fileName, content, int64(length), mime)
			So(err, ShouldBeNil)
		})
		Convey("Fail PutObject", func() {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(400)
			}))

			url := strings.TrimPrefix(server.URL, "http://")
			config := Config{
				AccessKeyID:     "x",
				Endpoint:        url,
				Region:          "x",
				SecretAccessKey: "x",
				BucketName:      "x",
				SSL:             false,
			}
			s3, err := New(config)
			So(err, ShouldBeNil)

			bucket := "string"
			directory := "string"
			fileName := "string.png"
			content := bytes.NewReader([]byte("asdf"))
			length := content.Len()
			mime := "image/png"
			err = s3.CreateFile(bucket, directory, fileName, content, int64(length), mime)
			So(err, ShouldNotBeNil)
		})
		Convey("Fail SetBucketPolicy", func() {
			i := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if i > 0 {
					w.WriteHeader(400)
				}
				fmt.Fprintln(w, "{}")
				i++

			}))

			url := strings.TrimPrefix(server.URL, "http://")
			config := Config{
				AccessKeyID:     "x",
				Endpoint:        url,
				Region:          "x",
				SecretAccessKey: "x",
				BucketName:      "x",
				SSL:             false,
			}
			s3, err := New(config)
			So(err, ShouldBeNil)

			bucket := "string"
			directory := "string"
			fileName := "string.png"
			content := bytes.NewReader([]byte("asdf"))
			length := content.Len()
			mime := "image/png"
			err = s3.CreateFile(bucket, directory, fileName, content, int64(length), mime)
			So(err, ShouldNotBeNil)
		})
	})
	Convey("GetS3Host", t, func() {
		endpoint := "localhost"
		config := Config{
			AccessKeyID:     "x",
			Endpoint:        endpoint,
			Region:          "x",
			SecretAccessKey: "x",
			BucketName:      "x",
			SSL:             false,
		}
		s3, err := New(config)
		So(err, ShouldBeNil)
		So(s3.GetS3Host(), ShouldEqual, endpoint)
	})

	Convey("BucketExists", t, func() {
		Convey("invalid response", func() {
			endpoint := "invalidhost-asdasd"
			config := Config{
				AccessKeyID:     "x",
				Endpoint:        endpoint,
				Region:          "x",
				SecretAccessKey: "x",
				BucketName:      "x",
				SSL:             false,
			}
			s3, err := New(config)
			So(err, ShouldBeNil)
			res, err := s3.BucketExists("somebucketname")
			So(res, ShouldBeFalse)
			So(err, ShouldNotBeNil)

		})
		Convey("Disabled S3", func() {
			s3 := helper{
				Enabled: false,
			}

			res, err := s3.BucketExists("x")
			So(err, ShouldBeNil)
			So(res, ShouldBeFalse)
		})
	})
}
