package internal

import "testing"

func TestValidateRPCSetCornerRadius(t *testing.T) {
	t.Run("accepts uniform corner radius", func(t *testing.T) {
		err := ValidateRPC("set_corner_radius", []string{"1:2"}, map[string]interface{}{
			"cornerRadius": 24.0,
		})
		if err != "" {
			t.Fatalf("expected no validation error, got %q", err)
		}
	})

	t.Run("accepts per corner radius", func(t *testing.T) {
		err := ValidateRPC("set_corner_radius", []string{"1:2"}, map[string]interface{}{
			"topLeftRadius":     24.0,
			"topRightRadius":    24.0,
			"bottomRightRadius": 8.0,
			"bottomLeftRadius":  8.0,
		})
		if err != "" {
			t.Fatalf("expected no validation error, got %q", err)
		}
	})

	t.Run("rejects missing radius params", func(t *testing.T) {
		err := ValidateRPC("set_corner_radius", []string{"1:2"}, map[string]interface{}{})
		want := "at least one corner radius property is required"
		if err != want {
			t.Fatalf("expected %q, got %q", want, err)
		}
	})

	t.Run("rejects negative radius", func(t *testing.T) {
		err := ValidateRPC("set_corner_radius", []string{"1:2"}, map[string]interface{}{
			"cornerRadius": -1.0,
		})
		want := "cornerRadius must be non-negative"
		if err != want {
			t.Fatalf("expected %q, got %q", want, err)
		}
	})
}

func TestValidateRPCSetFills(t *testing.T) {
	t.Run("accepts legacy solid fill", func(t *testing.T) {
		err := ValidateRPC("set_fills", []string{"1:2"}, map[string]interface{}{
			"color":   "#FF5733",
			"opacity": 0.8,
		})
		if err != "" {
			t.Fatalf("expected no validation error, got %q", err)
		}
	})

	t.Run("accepts gradient fills", func(t *testing.T) {
		err := ValidateRPC("set_fills", []string{"1:2"}, map[string]interface{}{
			"fills": []interface{}{
				map[string]interface{}{
					"type":  "GRADIENT_LINEAR",
					"angle": 135.0,
					"stops": []interface{}{
						map[string]interface{}{"position": 0.0, "color": "#0F172A"},
						map[string]interface{}{"position": 1.0, "color": "#F8FAFC", "opacity": 0.92},
					},
				},
			},
		})
		if err != "" {
			t.Fatalf("expected no validation error, got %q", err)
		}
	})

	t.Run("accepts empty fills array to clear fills", func(t *testing.T) {
		err := ValidateRPC("set_fills", []string{"1:2"}, map[string]interface{}{
			"fills": []interface{}{},
		})
		if err != "" {
			t.Fatalf("expected no validation error, got %q", err)
		}
	})

	t.Run("rejects missing fill payload", func(t *testing.T) {
		err := ValidateRPC("set_fills", []string{"1:2"}, map[string]interface{}{})
		want := "either fills or color is required"
		if err != want {
			t.Fatalf("expected %q, got %q", want, err)
		}
	})

	t.Run("rejects invalid gradient stop position", func(t *testing.T) {
		err := ValidateRPC("set_fills", []string{"1:2"}, map[string]interface{}{
			"fills": []interface{}{
				map[string]interface{}{
					"type": "GRADIENT_LINEAR",
					"stops": []interface{}{
						map[string]interface{}{"position": -0.1, "color": "#0F172A"},
						map[string]interface{}{"position": 1.0, "color": "#F8FAFC"},
					},
				},
			},
		})
		want := "fills[0].stops[0].position must be between 0 and 1"
		if err != want {
			t.Fatalf("expected %q, got %q", want, err)
		}
	})
}
