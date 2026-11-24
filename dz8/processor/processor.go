package processor

import "sync"

// ProcessFiles возвращает по результату на каждый переданный путь.
// Получатель читает канал до закрытия.
func ProcessFiles(paths []string, workers int) <-chan Result {
	if workers < 1 {
		workers = 1
	}
	jobs := make(chan string)
	results := make(chan Result)
	// TODO: распределите пути между workers обработчиками через канал jobs.
	// Для каждого пути вызовите analyzeFile и отправьте результат в results.
	// Определите, кто закроет каждый канал и как для этого использовать WaitGroup.
	//
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				result := analyzeFile(path)
				results <- result
			}
		}()
	}

	go func() {
		for _, path := range paths {
			jobs <- path
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}
