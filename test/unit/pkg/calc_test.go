package calc_test

import (
	"testing"

	"github.com/SobolevTim/finance_bot/internal/pkg/calc"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       float64
		wantErr    bool
	}{
		{
			name:       "simple",
			expression: "2+2",
			want:       4,
			wantErr:    false,
		},
		{
			name:       "simple with spaces",
			expression: " 2 + 2 ",
			want:       4,
			wantErr:    false,
		},
		{
			name:       "simple with brackets",
			expression: "(3-1)*2",
			want:       4,
			wantErr:    false,
		},
		{
			name:       "complex with brackets and spaces",
			expression: " ( 3 - 1 ) * 2 ",
			want:       4,
			wantErr:    false,
		},
		{
			name:       "simple with percent",
			expression: "100+10%",
			want:       110,
			wantErr:    false,
		},
		{
			name:       "simple with percent and spaces",
			expression: " 100 + 10 % ",
			want:       110,
			wantErr:    false,
		},
		{
			name:       "complex with percent and brackets",
			expression: "(100-10%)*2-50%",
			want:       90,
			wantErr:    false,
		},
		{
			name:       "complex with percent and brackets and spaces",
			expression: " ( 100 - 10 % ) * 2 - 50 % ",
			want:       90,
			wantErr:    false,
		},
		{
			name:       "simple pow",
			expression: "2^2",
			want:       4,
			wantErr:    false,
		},
		{
			name:       "simple pow with spaces",
			expression: " 2 ^ 2 ",
			want:       4,
			wantErr:    false,
		},
		{
			name:       "complex pow",
			expression: "2^2^2",
			want:       16,
			wantErr:    false,
		},
		{
			name:       "complex pow with persent",
			expression: "2^2^2+50%",
			want:       24,
			wantErr:    false,
		},
		{
			name:       "complex with milti percent",
			expression: "(100-10%)*2-50%",
			want:       90,
			wantErr:    false,
		},
		{
			name:       "error: division by zero",
			expression: "1/0",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "error: invalid expression",
			expression: "1/0+",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "error: invalid expression with brackets",
			expression: "(1/0+1)",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "error: invalid expression with percent",
			expression: "100+10%+",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "error: invalid expression with pow",
			expression: "2^2^",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "error: invalid expression with pow and percent",
			expression: "2^2^2+50%+",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "error: incorrect brackets",
			expression: "1+1))",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "error: incorrect brackets with spaces",
			expression: " 1 + 1 ) ) ",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "error: percentage with multiplication operator",
			expression: "100*10%",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "error: percentage without preceding operator",
			expression: "100%",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "unbalanced parentheses auto fix",
			expression: "((1+1)",
			want:       2,
			wantErr:    false,
		},
		{
			name:       "invalid symbol",
			expression: "1+1a",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "decimal numbers",
			expression: "1.1+1.1",
			want:       2.2,
			wantErr:    false,
		},
		{
			name:       "decimal numbers with percent",
			expression: "1.5+10.5%",
			want:       1.6575,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calc.Calculate(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate() test %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name string
		num  float64
		want string
	}{
		{
			name: "integer",
			num:  5.0,
			want: "5",
		},
		{
			name: "decimal with rounding",
			num:  3.1415926535,
			want: "3.14159",
		},
		{
			name: "no decimal part after rounding",
			num:  2.0000000001,
			want: "2",
		},
		{
			name: "rounding to five decimal places",
			num:  2.123456789,
			want: "2.12346",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.FormatNumber(tt.num)
			if got != tt.want {
				t.Errorf("FormatNumber(%v) = %v, want %v", tt.num, got, tt.want)
			}
		})
	}
}
