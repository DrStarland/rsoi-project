package handlers

// import (
// 	"sort"
// )

// type Users []User

// func (source Users) FindByQueryAndGetSlice(query string) Users {
// 	if query == "" {
// 		return source
// 	}
// 	newUsers := make(Users, 0)
// 	for _, user := range source {
// 		if user.containQuery(query) {
// 			newUsers = append(newUsers, user)
// 		}
// 	}
// 	return newUsers
// }

// func (source Users) DoOffset(offset int) Users {
// 	if offset <= 0 {
// 		return source
// 	}

// 	newUsers := make(Users, len(source)-offset)
// 	for i := range newUsers {
// 		newUsers[i] = source[i+offset]
// 	}
// 	return newUsers
// }

// func (source Users) CutToLimit(limit int) Users {
// 	if limit <= 0 {
// 		return make(Users, 0)
// 	}
// 	if limit >= len(source) {
// 		return source
// 	}

// 	return source[0:limit]
// }

// func (source Users) Sort(field string, direction int) Users {
// 	if direction == OrderByAsIs {
// 		return source
// 	}

// 	var f func(i, j int) bool

// 	switch direction {
// 	case OrderByAsc:
// 		switch field {
// 		case caseID:
// 			f = func(i, j int) bool { return source[i].ID < source[j].ID }
// 		case caseAge:
// 			f = func(i, j int) bool { return source[i].Age < source[j].Age }
// 		case caseName:
// 			f = func(i, j int) bool { return source[i].Name < source[j].Name }
// 		}
// 	case OrderByDesc:
// 		switch field {
// 		case caseID:
// 			f = func(i, j int) bool { return source[i].ID > source[j].ID }
// 		case caseAge:
// 			f = func(i, j int) bool { return source[i].Age > source[j].Age }
// 		case caseName:
// 			f = func(i, j int) bool { return source[i].Name > source[j].Name }
// 		}
// 	default:
// 		return source
// 	}

// 	sort.Slice(source, f)
// 	return source
// }
