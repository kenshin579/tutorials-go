package local

//
//var (
//	scheduleStore domain.ScheduleStore
//)
//
//func setup() {
//	c, _ := cronner.New()
//	scheduleStore = NewLocalScheduleStore(c)
//}
//
//func TestCreate_JotTypeHttpPost(t *testing.T) {
//	setup()
//	request := domain.ScheduleRequest{
//		JobDescription: "test-job",
//		JobType:        domain.JobTypeHttpPost,
//		Schedule:       "* * * * *",
//		JobRequest: domain.JobRequest{
//			Url: "",
//			Body: domain.MoveFileRequest{
//				FileName:    "test.txt",
//				Destination: "/test",
//			},
//		},
//	}
//
//	err := scheduleStore.Create(request)
//	assert.NoError(t, err)
//}
//
//func TestCreate_JotTypePrint(t *testing.T) {
//	//GIVEN
//	setup()
//	request := domain.ScheduleRequest{
//		JobDescription: "test-job",
//		JobType:        domain.JobTypePrint,
//		Schedule:       "* * * * * *",
//		JobRequest:     domain.JobRequest{},
//	}
//
//	//WHEN
//	err := scheduleStore.Create(request)
//
//	//THEN
//	assert.NoError(t, err)
//	list, _ := scheduleStore.List()
//	assert.Equal(t, 1, len(list))
//}
