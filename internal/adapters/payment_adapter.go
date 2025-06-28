package adapters

import (
	"context"
	"fmt"
	"snakers-bot/internal/interfaces"
	entity "snakers-bot/internal/usecases"
)

type SimplePaymentAdapter struct {
	companyName   string
	bankDetails   string // "СБЕР 2202 2002 2002 2002"
	cryptoWallets map[string]string
}

func NewSimplePaymentAdapter(companyName, bankDetails string, cryptoWallets map[string]string) interfaces.PaymentProvider {
	return &SimplePaymentAdapter{
		companyName:   companyName,
		bankDetails:   bankDetails,
		cryptoWallets: cryptoWallets,
	}
}

func (a *SimplePaymentAdapter) GeneratePaymentDetails(ctx context.Context, order *entity.Order) (string, error) {
	message := fmt.Sprintf(
		"💳 *Оплата заказа #%d*\n\n"+
			"Сумма: *%.2f руб.*\n\n"+
			"*Банковский перевод:*\n%s\n\n"+
			"*Криптовалюты:*\n",
		order.ID,
		order.TotalPrice(),
		a.bankDetails,
	)

	for currency, wallet := range a.cryptoWallets {
		message += fmt.Sprintf("%s: `%s`\n", currency, wallet)
	}

	message += "\nПосле оплаты пришлите чек менеджеру(проект тестовый) @Lev_arino"
	return message, nil
}
