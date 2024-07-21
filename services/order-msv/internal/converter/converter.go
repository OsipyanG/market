package converter

import (
	orderpb "github.com/OsipyanG/market/protos/order"
	"github.com/OsipyanG/market/services/order-msv/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetProtoOrder(modelOrder *model.Order) *orderpb.Order {
	items := GetProtoOrderItems(modelOrder.Items)

	return &orderpb.Order{
		OrderId:    modelOrder.ID.String(),
		CustomerId: modelOrder.CustomerID.String(),
		Status:     modelOrder.Status,
		Address:    modelOrder.Address,
		Items:      items,
		CreatedAt:  timestamppb.New(modelOrder.CreatedAt),
		UpdatedAt:  timestamppb.New(modelOrder.UpdatedAt),
	}
}

func GetProtoOrderItems(items []model.OrderItem) []*orderpb.OrderItem {
	protoItems := make([]*orderpb.OrderItem, 0, len(items))

	for _, item := range items {
		protoItems = append(protoItems, &orderpb.OrderItem{
			ProductId: item.ProductID.String(),
			Quantity:  int32(item.Quantity),
			Price:     uint64(item.Price),
		})
	}

	return protoItems
}

func GetProtoDelivery(modelDelivery model.Delivery) *orderpb.Delivery {
	return &orderpb.Delivery{
		OrderId:   modelDelivery.OrderID.String(),
		CourierId: modelDelivery.CourierID.String(),
		Status:    modelDelivery.Status,
	}
}
