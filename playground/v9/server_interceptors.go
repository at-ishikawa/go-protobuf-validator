package grpc_playground_validator

import (
    `context`

    `google.golang.org/grpc`
    `google.golang.org/grpc/codes`
    `google.golang.org/grpc/status`
    validatorv9 "gopkg.in/go-playground/validator.v9"
)

// UnaryServerInterceptor returns validator server interceptor for every request.
// If the request cannot be validated, the client gets the error with codes.FailedCondition
// If the request is nil, then this interceptor does nothing.
// If multiple languages are supported and should be different from each request,
// then the context which is generated by NewContextWithLocale should be passed before this interceptor.
func UnaryServerInterceptor(v *Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if req == nil {
			return handler(ctx, req)
		}

        st, err := v.ValidateGRPCRequest(ctx, req)
		if err != nil {
		    if inverr, ok := err.(*validatorv9.InvalidValidationError); ok {
		        return nil, status.Errorf(codes.FailedPrecondition, "request is not able to validate: %s", inverr.Error())
            }
			return nil, status.Errorf(codes.Internal, "invalid request for validator: %s", err.Error())
		}
		if st != nil {
			return nil, st.Err()
		}

		return handler(ctx, req)
	}
}
