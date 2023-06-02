package util

//
//import "fmt"
//
//type Queue []interface{}
//
//// IsEmpty - 큐가 비어있는지 확인하는 함수.
//func (q *Queue) IsEmpty() bool {
//	return len(*q) == 0
//}
//
//// Enqueue - 큐에 값을 추가하는 함수.
//func (q *Queue) Enqueue(data interface{}) {
//	*q = append(*q, data) // 큐 끝에 값을 추가함.
//}
//
//// Dequeue - 큐에 첫번째 요소를 반환하고 제거하는 함수.
//func (q *Queue) Dequeue() interface{} {
//	if q.IsEmpty() {
//		fmt.Println("queue is empty")
//		return nil
//	}
//	data := (*q)[0] // 큐에 첫번째 값을 가져옴.
//	*q = (*q)[1:]   // 큐에 첫번째 데이터를 제거함.
//	return data
//}
