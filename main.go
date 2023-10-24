package main

import (
	"fmt"
	"math"

	"github.com/devfans/golang/log"
)

func Mul(a, b int) int {
	return (a * b/1000 + 5) / 10
}

type TaxLevels [][3]int

func (tl TaxLevels) CalTax(taxSalary, accTaxSalary int) (tax int) {
	if taxSalary <= 0 { return }
	if accTaxSalary < 0 {
		accTaxSalary = 0
	}
	total := taxSalary + accTaxSalary
	accSal := 0
	accTax := 0
	for _, entry := range tl {
		if total < entry[0] {
			tax1 := Mul(total, entry[1]) - entry[2]
			tax2 := Mul((total - accSal), entry[1]) + accTax
			if tax1 != tax2 {
				panic(fmt.Errorf("cal tax failure %v != %v", tax1, tax2))
			} else {
				return tax1 - tl.CalTax(accTaxSalary, 0)
			}
		} else {
			accTax += Mul((entry[0] - accSal), entry[1])
			accSal = entry[0] 
		}
	}
	panic("overflow")
}


type PersonalMonthBill struct {
	Month int
	DutyFree int
	FullSalary, AccFullSalary int
	NetSalary, AccNetSalary int
	InsuranceFee, AccInsuranceFee int
	HouseFundingFee, AccHourseFundingFee int
	ExtraFees, AccExtraFees int
	TaxSalary, AccTaxSalary int
    Tax, AccTax int
}

func (b PersonalMonthBill) PrepareForMonth(month int) (bill PersonalMonthBill) {
	bill = PersonalMonthBill{
		Month: month,
		DutyFree: b.DutyFree,
		FullSalary: b.FullSalary,
		InsuranceFee: b.InsuranceFee,
		HouseFundingFee: b.HouseFundingFee,
		ExtraFees: b.ExtraFees,
		TaxSalary: b.FullSalary - b.InsuranceFee - b.HouseFundingFee - b.ExtraFees - b.DutyFree,
	}
	return
}

func (b PersonalMonthBill) PrepareForMonthWithAcc(month int, prev PersonalMonthBill) (bill PersonalMonthBill) {
	bill = b.PrepareForMonth(month)
	bill.AccFullSalary = prev.AccFullSalary + prev.FullSalary
	bill.AccNetSalary = prev.AccNetSalary + prev.NetSalary
	bill.AccInsuranceFee = prev.AccInsuranceFee + prev.InsuranceFee
	bill.AccHourseFundingFee = prev.AccHourseFundingFee + prev.HouseFundingFee
	bill.AccExtraFees = prev.AccExtraFees + prev.ExtraFees
	bill.AccTaxSalary = prev.AccTaxSalary + prev.TaxSalary
	bill.AccTax = prev.AccTax + prev.Tax
	return
}

type IPersonalTaxCalucator interface {
	CalForMonth(month int) PersonalMonthBill
}

type SHPersonalTaxCalculator struct {
	StartMonth *PersonalMonthBill
	MonthlyBill PersonalMonthBill
	TaxLevels TaxLevels
	InitMonth *int
	ReadyBills []*PersonalMonthBill
}

func (c *SHPersonalTaxCalculator) CalForMonth(month int) (b PersonalMonthBill) {
	if month > 12 || month < 1 {
		panic(fmt.Errorf("invalid month %d", month))
	}

	if len(c.ReadyBills) >= month && c.ReadyBills[month-1] != nil {
		return *c.ReadyBills[month-1]
	}

	if c.StartMonth != nil && month == c.StartMonth.Month {
		b = *c.StartMonth
		if b.Tax == 0 {
			b.TaxSalary = b.FullSalary - b.InsuranceFee - b.HouseFundingFee - b.ExtraFees - b.DutyFree
			b.Tax = c.TaxLevels.CalTax(b.TaxSalary, b.AccTaxSalary)
			b.NetSalary = b.TaxSalary - b.Tax + b.DutyFree
		}
		if len(c.ReadyBills) >= month {
			cp := b
			c.ReadyBills[month-1] = &cp
		}
		return
	}
	initMonth := 1
	if c.InitMonth != nil {
		initMonth = *c.InitMonth
	}
	if month == initMonth {
		b = c.MonthlyBill.PrepareForMonth(month)
	} else {
		prev := c.CalForMonth(month-1)
		b = c.MonthlyBill.PrepareForMonthWithAcc(month, prev)
	}
	b.Tax = c.TaxLevels.CalTax(b.TaxSalary, b.AccTaxSalary)
	b.NetSalary = b.TaxSalary - b.Tax + b.DutyFree
	if len(c.ReadyBills) >= month {
		cp := b
		c.ReadyBills[month-1] = &cp
	}
	return
}

func main() {
	monthlyBill := PersonalMonthBill{
		DutyFree: 500000,
		FullSalary: 3660000,
		InsuranceFee: 383765,
		HouseFundingFee: 438500,
	}
	startMonth := monthlyBill
	startMonth.Month = 6
	startMonth.InsuranceFee = 358974
	startMonth.HouseFundingFee = 353600

	cal := SHPersonalTaxCalculator {
		StartMonth: &startMonth,
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
		log.Info("MonthlyBill", "month", b.Month, "full_salary", b.FullSalary, "net_salary", b.NetSalary, "tax", b.Tax, "tax_salary", b.TaxSalary)
		log.Json(log.INFO, b)
	}
}