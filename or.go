package done

// Or объединяет несколько каналов done в один общий канал,
// который закрывается, как только закроется любой из переданных каналов.
// Полезно для агрегирования сигналов завершения из нескольких источников,
// позволяя одному слушателю ждать закрытия любого из каналов.
// Функция рекурсивно объединяет каналы с помощью select,
// возвращая канал, который закрывается однократно при первом сигнале.
// Параметры:
//
//	channels - один или более каналов типа <-chan interface{}
//
// Возвращаемое значение:
//
//	<-chan interface{} - канал, закрывающийся при закрытии любого из входных каналов.
func Or(channels ...<-chan interface{}) <-chan interface{} {
	mergedChan := make(chan interface{})

	switch len(channels) {
	case 0:
		close(mergedChan)
		return mergedChan
	case 1:
		return channels[0]
	case 2:
		go func() {
			defer close(mergedChan)

			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		}()
		return mergedChan
	default:
		go func() {
			defer close(mergedChan)

			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-Or(channels[2:]...):
			}
		}()
		return mergedChan
	}
}
