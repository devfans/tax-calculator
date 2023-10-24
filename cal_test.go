package main

import (
	"math"
	"testing"

	"github.com/devfans/golang/log"
)

func TestMain(t *testing.T) {
	monthlyBill := PersonalMonthBill{
		DutyFree:        500000,
		FullSalary:      3660000,
		InsuranceFee:    383765,
		HouseFundingFee: 438500,
	}
	startMonth := monthlyBill
	startMonth.Month = 6
	startMonth.InsuranceFee = 358974
	startMonth.HouseFundingFee = 353600

	cal := SHPersonalTaxCalculator{
		StartMonth:  &startMonth,
		MonthlyBill: monthlyBill,
		TaxLevels: TaxLevels{
			{3600000, 300, 0},
			{14400000, 1000, 252000},
			{30000000, 2000, 1692000},
			{42000000, 2500, 3192000},
			{66000000, 3000, 5292000},
			{96000000, 3500, 8592000},
			{math.MaxInt, 4500, 18192000},
		},
		ReadyBills: make([]*PersonalMonthBill, 12),
	}
	cal.CalForMonth(12)
	for _, b := range cal.ReadyBills {
		if b == nil {
			continue
		}
		log.Info("MonthlyBill", "month", b.Month, "full_salary", b.FullSalary, "net_salary", b.NetSalary, "tax", b.Tax, "acc_tax_salary", b.AccTaxSalary)
		log.Json(log.INFO, b)
	}
}

func TestMainFull(t *testing.T) {
	monthlyBill := PersonalMonthBill{
		DutyFree:        500000,
		FullSalary:      3660000,
		InsuranceFee:    383765,
		HouseFundingFee: 438500,
	}

	cal := SHPersonalTaxCalculator{
		MonthlyBill: monthlyBill,
		TaxLevels: TaxLevels{
			{3600000, 300, 0},
			{14400000, 1000, 252000},
			{30000000, 2000, 1692000},
			{42000000, 2500, 3192000},
			{66000000, 3000, 5292000},
			{96000000, 3500, 8592000},
			{math.MaxInt, 4500, 18192000},
		},
		ReadyBills: make([]*PersonalMonthBill, 12),
	}
	cal.CalForMonth(12)
	for _, b := range cal.ReadyBills {
		if b == nil {
			continue
		}
		log.Info("MonthlyBill", "month", b.Month, "full_salary", b.FullSalary, "net_salary", b.NetSalary, "tax", b.Tax, "acc_tax_salary", b.AccTaxSalary)
		// log.Json(log.INFO, b)
	}
}
