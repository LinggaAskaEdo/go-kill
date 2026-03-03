package grpc

// func (g *Grpc) CreateAuthUser(ctx context.Context, req *authpb.CreateAuthUserRequest) (*authpb.CreateAuthUserResponse, error) {
// 	resp, err := g.authClient.CreateAuthUser(ctx, &authpb.CreateAuthUserRequest{
// 		Email:    req.Email,
// 		Password: req.Password,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }

// func (g *Grpc) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
// 	resp, err := g.authClient.ValidateToken(ctx, &authpb.ValidateTokenRequest{
// 		Token: req.Token,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }
