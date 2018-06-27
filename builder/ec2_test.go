package builder

import "testing"

func TestIncrementID(t *testing.T) {
	instances := map[string]InstanceInfo{
		"test": InstanceInfo{InstanceID: "i-123456"},
	}
	dup_instance := IncrementID("test", instances)
	if dup_instance != "test-2" {
		t.Errorf("Duplicate instnace id should append unique num, got: %s", dup_instance)
	}
}

func TestIncrementID2(t *testing.T) {
	instances := map[string]InstanceInfo{
		"test-2": InstanceInfo{InstanceID: "i-123456"},
	}
	dup_instance := IncrementID("test-2", instances)
	if dup_instance != "test-3" {
		t.Errorf("Duplicate instnace id should append unique num, got: %s", dup_instance)
	}
}
