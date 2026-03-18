package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"
	"github.com/GoogleCloudPlatform/microservices-demo/src/frontend/money"
	"golang.org/x/sync/errgroup"
)

// Simplified sequential viewCart logic
func benchmarkViewCartLogicSequential(ctx context.Context, cart []*pb.CartItem) error {
	type cartItemView struct {
		Item     *pb.Product
		Quantity int32
		Price    *pb.Money
	}
	items := make([]cartItemView, len(cart))
	totalPrice := pb.Money{CurrencyCode: "USD"}
	for i, item := range cart {
		p, err := mockGetProduct(ctx, item.GetProductId())
		if err != nil {
			return err
		}
		price, err := mockConvertCurrency(ctx, p.GetPriceUsd(), "USD")
		if err != nil {
			return err
		}

		multPrice := money.MultiplySlow(*price, uint32(item.GetQuantity()))
		items[i] = cartItemView{
			Item:     p,
			Quantity: item.GetQuantity(),
			Price:    &multPrice}
		totalPrice = money.Must(money.Sum(totalPrice, multPrice))
	}
	return nil
}

// Simplified concurrent viewCart logic
func benchmarkViewCartLogicConcurrent(ctx context.Context, cart []*pb.CartItem) error {
	type cartItemView struct {
		Item     *pb.Product
		Quantity int32
		Price    *pb.Money
	}
	items := make([]cartItemView, len(cart))
	g, ctx := errgroup.WithContext(ctx)
	for i, item := range cart {
		i, item := i, item
		g.Go(func() error {
			p, err := mockGetProduct(ctx, item.GetProductId())
			if err != nil {
				return err
			}
			price, err := mockConvertCurrency(ctx, p.GetPriceUsd(), "USD")
			if err != nil {
				return err
			}

			multPrice := money.MultiplySlow(*price, uint32(item.GetQuantity()))
			items[i] = cartItemView{
				Item:     p,
				Quantity: item.GetQuantity(),
				Price:    &multPrice}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	totalPrice := pb.Money{CurrencyCode: "USD"}
	for _, item := range items {
		totalPrice = money.Must(money.Sum(totalPrice, *item.Price))
	}
	return nil
}

func mockGetProduct(ctx context.Context, id string) (*pb.Product, error) {
	time.Sleep(10 * time.Millisecond) // Simulate latency
	return &pb.Product{
		Id:   id,
		Name: "Mock Product " + id,
		PriceUsd: &pb.Money{
			CurrencyCode: "USD",
			Units:        10,
			Nanos:        0,
		},
	}, nil
}

func mockConvertCurrency(ctx context.Context, from *pb.Money, to string) (*pb.Money, error) {
	time.Sleep(5 * time.Millisecond) // Simulate latency
	return &pb.Money{
		CurrencyCode: to,
		Units:        from.Units,
		Nanos:        from.Nanos,
	}, nil
}

func TestBenchmarkComparison(t *testing.T) {
	ctx := context.Background()
	cart := make([]*pb.CartItem, 10)
	for i := 0; i < 10; i++ {
		cart[i] = &pb.CartItem{
			ProductId: fmt.Sprintf("prod-%d", i),
			Quantity:  int32(i + 1),
		}
	}

	startSeq := time.Now()
	err := benchmarkViewCartLogicSequential(ctx, cart)
	if err != nil {
		t.Fatal(err)
	}
	elapsedSeq := time.Since(startSeq)
	fmt.Printf("Sequential loop for 10 items took %v\n", elapsedSeq)

	startConc := time.Now()
	err = benchmarkViewCartLogicConcurrent(ctx, cart)
	if err != nil {
		t.Fatal(err)
	}
	elapsedConc := time.Since(startConc)
	fmt.Printf("Concurrent loop for 10 items took %v\n", elapsedConc)
}
