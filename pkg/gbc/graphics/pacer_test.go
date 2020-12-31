package graphics

import (
	"time"
)

func (s *unitTestSuite) TestPacer() {
	s.Run("some sleep required is done", func() {
		s.tctx.mockClock.EXPECT().Since(s.tctx.originalTime).Return(1400000 * time.Nanosecond)
		s.tctx.mockClock.EXPECT().Sleep(15200000 * time.Nanosecond)
		s.tctx.mockClock.EXPECT().Now().Return(s.tctx.originalTime.Add(17 * time.Millisecond))

		s.tctx.gpu.pacer.Wait()
	})

	s.Run("no pacing necessary makes sleep skipped", func() {
		s.tctx.mockClock.EXPECT().Since(s.tctx.originalTime).Return(time.Second)
		s.tctx.mockClock.EXPECT().Now().Return(s.tctx.originalTime.Add(time.Second))

		s.tctx.gpu.pacer.Wait()
	})
}
