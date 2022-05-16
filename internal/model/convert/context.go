package convert

//func Context(eCtx echo.Context) context.Context {
//	ctx := context.Background()
//
//	userID := eCtx.Request().Header.Get(constants.HeaderKeyUserID)
//	if len(userID) != 0 {
//		ctx = context.WithValue(ctx, constants.CtxKeyUserID{}, userID)
//	}
//
//	xRequestID := eCtx.Request().Header.Get(constants.HeaderKeyRequestID)
//	if len(userID) != 0 {
//		ctx = context.WithValue(ctx, constants.CtxKeyXRequestID{}, xRequestID)
//	}
//
//	return ctx
//}
