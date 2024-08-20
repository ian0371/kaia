package types

type Partition[T any] struct {
	boundary uint
	item     T
}

type PartitionList[T any] []Partition[T]

func (iv *PartitionList[T]) AddRecord(newBoundary uint, newValue T) {
	if len(*iv) == 0 {
		*iv = append(*iv, Partition[T]{newBoundary, newValue})
		return
	}

	for i, interval := range *iv {
		if interval.boundary == newBoundary {
			// Update existing value if boundary already exists
			(*iv)[i].item = newValue
			return
		}
		if interval.boundary > newBoundary {
			// Insert new interval at the correct position
			*iv = append(*iv, Partition[T]{})
			copy((*iv)[i+1:], (*iv)[i:])
			(*iv)[i] = Partition[T]{newBoundary, newValue}
			return
		}
	}

	// If we've reached here, add the new interval at the end
	*iv = append(*iv, Partition[T]{newBoundary, newValue})
}

func (iv *PartitionList[T]) RemoveRecord(i uint) {
	(*iv) = append((*iv)[:i], (*iv)[i+1:]...)
}

func (iv *PartitionList[T]) GetItem(height uint) T {
	index := 0
	for i := range *iv {
		if i == len(*iv)-1 {
			index = i
		} else if (*iv)[i].boundary <= height && height < (*iv)[i+1].boundary {
			index = i
			break
		}
	}

	return (*iv)[index].item
}
