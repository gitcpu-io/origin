package services

import (
	"context"
	"testing"

	"github.com/gitcpu-io/origin/models/ioparams"
	"github.com/gitcpu-io/origin/models/weather"
	"github.com/gitcpu-io/origin/services/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = Describe("Weather", func() {
	var s Weatherer
	var city string
	var svcobj *svc
	var req *ioparams.WeatherRequest
	Describe("List", func() {
		BeforeEach(func() {
			s = NewWeather()
			gomega.Expect(s).NotTo(gomega.BeNil())
			city = "深圳市"

			svcobj = &svc{}
			req = &ioparams.WeatherRequest{
				Query: "深圳市",
			}
		})
		Context("begin", func() {
			It("first test", func() {
				ctrl := gomock.NewController(GinkgoT())
				client := mocks.NewMockWeatherer(ctrl)
				var result *weather.Weather
				client.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(result, nil)

				_, err := s.List(context.TODO(), city)
				gomega.Expect(err).To(gomega.BeNil())

			})
		})
		Context("request", func() {
			It("htt test", func() {

				_, err := svcobj.dealRequestWeather(context.TODO(), req)
				gomega.Expect(err).To(gomega.BeNil())

			})
		})

	})

})

func BenchmarkDealRequestWeather(b *testing.B) {
	s := &svc{}
	req := &ioparams.WeatherRequest{
		Query: "深圳市",
	}
	for i := 0; i < b.N; i++ {
		s.dealRequestWeather(context.TODO(), req)
	}
}
