package main

// func (app *application) recoverPanic(next http.Handler) http.Handler {

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)
// 	{
// 		defer func(){
// 			err:=recover()
// 			if err!=nil {
// 				w.Header().Set("Connection","close")
// 				app.serverErrorResponse(w,r,fmt.Errorf("%s",err))
// 			}
// 		}()
// 	}
// 	){

// 	}
// }
