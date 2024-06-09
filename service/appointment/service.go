package appointment

import (
	"cmp"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/yageunpro/owl-backend-go/store"
	"github.com/yageunpro/owl-backend-go/store/appointment"
	"slices"
	"sort"
	"strconv"
	"time"
)

type Service interface {
	Add(ctx context.Context, arg AddParam) (*resInfo, error)
	List(ctx context.Context, arg ListParam) (*resInfoList, error)
	Info(ctx context.Context, appointmentId uuid.UUID) (*resInfo, error)
	Edit(ctx context.Context, arg EditParam) (*resInfo, error)
	Delete(ctx context.Context, appointmentId, userId uuid.UUID) error
	Share(ctx context.Context) error
	Join(ctx context.Context, appointmentId, userId uuid.UUID) (*resInfo, error)
	JoinNonmember(ctx context.Context, arg JoinNonmemberParam) (*resInfo, error)
	RecommendTime(ctx context.Context, appointmentId, userId uuid.UUID) ([]time.Time, error)
	Confirm(ctx context.Context, appointmentId, userId uuid.UUID, confirmTime time.Time) error
}
type service struct {
	store *store.Store
}

func New(sto *store.Store) Service {
	return &service{store: sto}
}

func (s *service) Add(ctx context.Context, arg AddParam) (*resInfo, error) {
	id := uuid.Must(uuid.NewV7())
	err := s.store.Appointment.Create(ctx, appointment.CreateAppointmentParam{
		Id:           id,
		UserId:       arg.UserId,
		Title:        arg.Title,
		Description:  arg.Description,
		LocationId:   arg.LocationId,
		CategoryList: arg.CategoryList,
		Deadline:     arg.Deadline,
	})
	if err != nil {
		return nil, err
	}
	return s.Info(ctx, id)
}

func (s *service) List(ctx context.Context, arg ListParam) (*resInfoList, error) {
	var offset int
	if arg.PageToken == nil || *arg.PageToken == "" {
		offset = 0
	} else {
		i, err := strconv.ParseInt(*arg.PageToken, 10, 32)
		if err != nil {
			offset = 0
		} else {
			offset = int(i)
		}
	}

	out, err := s.store.Appointment.List(ctx, appointment.Status(arg.Status), arg.UserId, offset, arg.Limit)
	if err != nil {
		return nil, err
	}

	offset += len(out)

	res := resInfoList{
		InfoList:  make([]absInfo, len(out)),
		NextToken: strconv.Itoa(offset),
	}

	for i := range out {
		var location *resLocation
		if out[i].LocationId.Valid {
			l, err := s.store.Location.GetLocation(ctx, out[i].LocationId.UUID)
			if err != nil {
				location = nil
			} else {
				location = &resLocation{
					Id:       l.Id,
					Title:    l.Title,
					Address:  l.Address,
					Category: l.Category,
					Position: [2]int{l.MapX, l.MapY},
				}
			}
		}

		res.InfoList[i] = absInfo{
			Id:          out[i].Id,
			OrganizerId: out[i].OrganizerId,
			Title:       out[i].Title,
			Location:    location,
			Status:      string(out[i].Status),
			ConfirmTime: out[i].ConfirmTime,
			HeadCount:   out[i].HeadCount,
		}
	}

	return &res, nil
}

func (s *service) Info(ctx context.Context, appointmentId uuid.UUID) (*resInfo, error) {
	out, err := s.store.Appointment.Get(ctx, appointmentId)
	if err != nil {
		return nil, err
	}

	participants := make([]resParticipant, len(out.Participants))
	for i := range out.Participants {
		participants[i].UserId = out.Participants[i].UserId
		participants[i].Username = out.Participants[i].Username
	}

	var location *resLocation
	if out.LocationId.Valid {
		l, err := s.store.Location.GetLocation(ctx, out.LocationId.UUID)
		if err != nil {
			location = nil
		} else {
			location = &resLocation{
				Id:       l.Id,
				Title:    l.Title,
				Address:  l.Address,
				Category: l.Category,
				Position: [2]int{l.MapX, l.MapY},
			}
		}
	}

	return &resInfo{
		Id:           out.Id,
		OrganizerId:  out.OrganizerId,
		Status:       string(out.Status),
		Title:        out.Title,
		Description:  out.Description,
		Deadline:     out.Deadline,
		ConfirmTime:  out.ConfirmTime,
		Location:     location,
		CategoryList: out.CategoryList,
		Participants: participants,
	}, nil
}

func (s *service) Edit(ctx context.Context, arg EditParam) (*resInfo, error) {
	locationId := uuid.NullUUID{Valid: false}
	if arg.LocationId != nil {
		locationId.UUID = *arg.LocationId
		locationId.Valid = true
	}
	err := s.store.Appointment.Update(ctx, appointment.UpdateAppointmentParam{
		Id:           arg.Id,
		UserId:       arg.UserId,
		Title:        arg.Title,
		Description:  arg.Description,
		LocationId:   locationId,
		CategoryList: arg.CategoryList,
	})
	if err != nil {
		return nil, err
	}

	return s.Info(ctx, arg.Id)
}

func (s *service) Delete(ctx context.Context, appointmentId, userId uuid.UUID) error {
	return s.store.Appointment.Delete(ctx, appointmentId, userId)
}

func (s *service) Share(ctx context.Context) error {
	return nil
}

func (s *service) Join(ctx context.Context, appointmentId, userId uuid.UUID) (*resInfo, error) {
	err := s.store.Appointment.AddParticipant(ctx, appointmentId, userId)
	if err != nil {
		return nil, err
	}

	return s.Info(ctx, appointmentId)
}

func (s *service) JoinNonmember(ctx context.Context, arg JoinNonmemberParam) (*resInfo, error) {
	param := appointment.AddParticipantNonMemberParam{
		Id:           arg.Id,
		Username:     arg.Username,
		ScheduleList: nil,
	}
	if arg.JoinSchedule == nil {
		param.ScheduleList = make([]appointment.Schedule, 0)
	} else {
		param.ScheduleList = make([]appointment.Schedule, len(arg.JoinSchedule))
		for i := range arg.JoinSchedule {
			param.ScheduleList[i] = appointment.Schedule{
				Title:     arg.JoinSchedule[i].Title,
				StartTime: arg.JoinSchedule[i].StartTime,
				EndTime:   arg.JoinSchedule[i].EndTime,
			}
		}
	}

	err := s.store.Appointment.AddParticipantNonMember(ctx, param)
	if err != nil {
		return nil, err
	}

	return s.Info(ctx, param.Id)
}

func (s *service) RecommendTime(ctx context.Context, appointmentId, userId uuid.UUID) ([]time.Time, error) {
	ap, err := s.store.Appointment.Get(ctx, appointmentId)
	if err != nil {
		return nil, err
	}

	if ap.OrganizerId != userId {
		return nil, errors.New("only organizer can get recommendations")
	}

	participantsIds := make([]uuid.UUID, 0)
	for i := range ap.Participants {
		participantsIds = append(participantsIds, ap.Participants[i].UserId)
	}

	if time.Now().UTC().After(ap.Deadline) {
		return nil, errors.New("appointment deadline exceeded")
	}

	res, err := s.internalRecommendTimes(ctx, participantsIds, time.Now().UTC(), ap.Deadline)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) Confirm(ctx context.Context, appointmentId, userId uuid.UUID, confirmTime time.Time) error {
	return s.store.Appointment.Confirm(ctx, appointmentId, userId, confirmTime)
}

func (s *service) internalRecommendTimes(ctx context.Context, userIds []uuid.UUID, startTime, endTime time.Time) ([]time.Time, error) {
	startTime = startTime.Truncate(time.Hour).Add(time.Hour).UTC()
	endTime = endTime.Truncate(time.Hour).Add(time.Hour).UTC()
	duration := time.Minute * 30

	out, err := s.store.Calendar.GetAllSchedule(ctx, userIds, startTime, endTime)
	if err != nil {
		return nil, err
	}

	krTimeZone := time.FixedZone("KR", 9*60*60)

	slots := make([]TimeSlot, 0)
	for startTime.Before(endTime) {
		if startTime.In(krTimeZone).Hour() < 10 || startTime.In(krTimeZone).Hour() > 19 {
			startTime = startTime.Add(duration)
			continue
		}

		slots = append(slots, TimeSlot{
			StartTime:       startTime,
			EndTime:         startTime.Add(duration),
			OverlapCount:    0,
			PreferenceScore: 0,
		})
		startTime = startTime.Add(duration)
	}

	for _, userSchedule := range out {
		userPeriod := TimePeriod{
			StartTime: userSchedule.StartTime,
			EndTime:   userSchedule.EndTime,
		}
		for i := range slots {
			slotPeriod := TimePeriod{
				StartTime: slots[i].StartTime,
				EndTime:   slots[i].EndTime,
			}
			if doTimePeriodsOverlap(userPeriod, slotPeriod) {
				slots[i].OverlapCount += 1
			}
		}
	}

	sort.Slice(slots, func(i, j int) bool {
		if slots[i].OverlapCount == slots[j].OverlapCount {
			return slots[i].StartTime.Before(slots[j].StartTime)
		}
		return slots[i].OverlapCount < slots[j].OverlapCount
	})

	for i := range slots {
		if i < len(slots)-1 {
			if slots[i].EndTime == slots[i+1].StartTime {
				slots[i].PreferenceScore += 1
			}
		}
	}

	mergeIndex := 0
	mergedSlots := make([]TimeSlot, 0)
	mergedSlots = append(mergedSlots, slots[0])

	for i := range slots {
		if i == 0 {
			continue
		}

		if i < len(slots)-1 {
			if mergedSlots[mergeIndex].EndTime == slots[i].StartTime && mergedSlots[mergeIndex].OverlapCount == slots[i].OverlapCount {
				mergedSlots[mergeIndex].EndTime = slots[i].EndTime
				mergedSlots[mergeIndex].PreferenceScore += slots[i].PreferenceScore
			} else {
				mergedSlots = append(mergedSlots, slots[i])
				mergeIndex += 1
			}
		}
	}

	keys := make([]int, 0)
	sortOverlap := make(map[int][]TimeSlot)
	for _, slot := range mergedSlots {
		_, ok := sortOverlap[slot.OverlapCount]
		if !ok {
			keys = append(keys, slot.OverlapCount)
			sortOverlap[slot.OverlapCount] = []TimeSlot{slot}
			continue
		}
		sortOverlap[slot.OverlapCount] = append(sortOverlap[slot.OverlapCount], slot)
	}

	slices.Sort(keys)

	amTimes := make([]TimeSlot, 0)
	pmTimes := make([]TimeSlot, 0)
	eveningTimes := make([]TimeSlot, 0)
	for key := range keys {
		slices.SortFunc(sortOverlap[key], func(a, b TimeSlot) int {
			return cmp.Compare(b.PreferenceScore, a.PreferenceScore)
		})

		for _, t := range sortOverlap[key] {
			if t.StartTime.In(krTimeZone).Hour() < 13 {
				amTimes = append(amTimes, t)
				continue
			} else if t.StartTime.In(krTimeZone).Hour() < 17 {
				pmTimes = append(pmTimes, t)
				continue
			} else {
				eveningTimes = append(eveningTimes, t)
			}
		}
	}

	resSlots := make([]TimeSlot, 0)
	done := 0
	resIndex := 0
	for {
		if done == 3 || len(resSlots) >= 6 {
			break
		}

		if resIndex == len(amTimes)-1 {
			done += 1
		} else if resIndex < len(amTimes) {
			resSlots = append(resSlots, amTimes[resIndex])
		}

		if resIndex == len(pmTimes)-1 {
			done += 1
		} else if resIndex < len(pmTimes) {
			resSlots = append(resSlots, pmTimes[resIndex])
		}

		if resIndex == len(resSlots)-1 {
			done += 1
		} else if resIndex < len(eveningTimes) {
			resSlots = append(resSlots, eveningTimes[resIndex])
		}

		resIndex += 1
	}

	res := make([]time.Time, len(resSlots))

	sort.Slice(resSlots, func(i, j int) bool {
		if resSlots[i].OverlapCount == resSlots[j].OverlapCount {
			return resSlots[i].PreferenceScore > resSlots[j].PreferenceScore
		}
		return resSlots[i].OverlapCount < resSlots[j].OverlapCount
	})

	for i := range resSlots {
		res[i] = resSlots[i].StartTime.In(krTimeZone)
	}

	return res, nil
}
