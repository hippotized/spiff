package dynaml

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hippotized/spiff/yaml"
)

func newNetworkFakeBinding(subnets yaml.Node, instances interface{}) Binding {
	return FakeBinding{
		FoundReferences: map[string]yaml.Node{
			"name":      node("cf1"),
			"instances": node(instances),
		},
		FoundFromRoot: map[string]yaml.Node{
			"":                     node("dummy"),
			"networks":             node("dummy"),
			"networks.cf1":         node("dummy"),
			"networks.cf1.subnets": subnets,
		},
	}
}

var _ = Describe("calls", func() {
	Describe("CIDR functions", func() {
		It("determines minimal IP", func() {
			expr := CallExpr{
				Name: "min_ip",
				Arguments: []Expression{
					StringExpr{"192.168.0.1/24"},
				},
			}

			Expect(expr).To(
				EvaluateAs(
					"192.168.0.0",
					FakeBinding{},
				),
			)
		})

		It("determines maximal IP", func() {
			expr := CallExpr{
				Name: "max_ip",
				Arguments: []Expression{
					StringExpr{"192.168.0.1/24"},
				},
			}

			Expect(expr).To(
				EvaluateAs(
					"192.168.0.255",
					FakeBinding{},
				),
			)
		})
	})

	Describe("join(\", \"...)", func() {
		expr := CallExpr{
			Name: "join",
			Arguments: []Expression{
				StringExpr{", "},
				ReferenceExpr{[]string{"alice"}},
				ReferenceExpr{[]string{"bob"}},
			},
		}

		It("joins string values ", func() {
			binding := FakeBinding{
				FoundReferences: map[string]yaml.Node{
					"alice": node("alice"),
					"bob":   node("bob"),
				},
			}

			Expect(expr).To(
				EvaluateAs(
					"alice, bob",
					binding,
				),
			)
		})

		It("joins int values ", func() {
			binding := FakeBinding{
				FoundReferences: map[string]yaml.Node{
					"alice": node(10),
					"bob":   node(20),
				},
			}

			Expect(expr).To(
				EvaluateAs(
					"10, 20",
					binding,
				),
			)
		})

		It("joins list entries ", func() {
			list := parseYAML(`
  - foo
  - bar
`)

			binding := FakeBinding{
				FoundReferences: map[string]yaml.Node{
					"alice": list,
					"bob":   node(20),
				},
			}

			Expect(expr).To(
				EvaluateAs(
					"foo, bar, 20",
					binding,
				),
			)
		})

		It("joins nothing", func() {
			expr := CallExpr{
				Name: "join",
				Arguments: []Expression{
					StringExpr{", "},
				},
			}

			Expect(expr).To(
				EvaluateAs(
					"",
					nil,
				),
			)
		})

		It("fails for missing args", func() {
			expr := CallExpr{
				Name:      "join",
				Arguments: []Expression{},
			}

			Expect(expr).To(FailToEvaluate(nil))
		})

		It("fails for wrong separator type", func() {
			expr := CallExpr{
				Name: "join",
				Arguments: []Expression{
					ListExpr{[]Expression{IntegerExpr{0}}},
				},
			}

			Expect(expr).To(FailToEvaluate(nil))
		})
	})

	Describe("static_ips(ips...)", func() {
		expr := CallExpr{
			Name: "static_ips",
			Arguments: []Expression{
				IntegerExpr{0},
				IntegerExpr{4},
			},
		}

		It("returns a set of ips from the given network's subnets", func() {
			subnets := parseYAML(`
- static:
    - 10.10.16.10
- static:
    - 10.10.16.11 - 10.10.16.254
`)
			binding := newNetworkFakeBinding(subnets, 2)

			Expect(expr).To(
				EvaluateAs(
					[]yaml.Node{node("10.10.16.10"), node("10.10.16.14")},
					binding,
				),
			)
		})

		It("limits the IPs to the number of instances", func() {
			subnets := parseYAML(`
- static:
    - 10.10.16.10 - 10.10.16.254
`)

			binding := newNetworkFakeBinding(subnets, 1)

			Expect(expr).To(
				EvaluateAs(
					[]yaml.Node{node("10.10.16.10")},
					binding,
				),
			)
		})

		Context("when the instance count is dynamic", func() {
			It("fails", func() {
				subnets := parseYAML(`
- static:
    - 10.10.16.10 - 10.10.16.254
`)

				binding := newNetworkFakeBinding(subnets, MergeExpr{})

				Expect(expr).To(FailToEvaluate(binding))
			})
		})

		Context("when there are not enough IPs for the number of instances", func() {
			It("fails", func() {
				subnets := parseYAML(`
- static:
    - 10.10.16.10 - 10.10.16.32
`)

				binding := newNetworkFakeBinding(subnets, 42)

				Expect(expr).To(FailToEvaluate(binding))
			})
		})

		Context("when there are singular static IPs listed", func() {
			It("includes them in the pool", func() {
				subnets := parseYAML(`
- static:
    - 10.10.16.10 - 10.10.16.32
    - 10.10.16.33
    - 10.10.16.34
`)

				expr := CallExpr{
					Name: "static_ips",
					Arguments: []Expression{
						IntegerExpr{0},
						IntegerExpr{4},
						IntegerExpr{23},
					},
				}

				binding := newNetworkFakeBinding(subnets, 3)

				Expect(expr).To(
					EvaluateAs(
						[]yaml.Node{node("10.10.16.10"), node("10.10.16.14"), node("10.10.16.33")},
						binding,
					),
				)
			})
		})
	})
})
