// the response to exchange a token key for a token
export type ExchangeTokenResponse = {
    // the recipes
    token: string | undefined;
};

// the request to exchange a token key for a token
export type ExchangeTokenRequest = {
    // the recipe to update
    tokenKey: string | undefined;
    // the fields to update
};

// the response to check if a token is valid
export type CheckTokenResponse = {
    // the user id on the token
    userId: number | undefined;
};

// the request to check if a token is valid
export type CheckTokenRequest = {}

// the auth service
export interface AuthService {
    // exchange a token key for a token
    ExchangeToken(request: ExchangeTokenRequest): Promise<ExchangeTokenResponse>;
    // check if a token is valid
    CheckToken(request: CheckTokenRequest): Promise<CheckTokenResponse>;
  }
  
type RequestType = {
    path: string;
    method: string;
    body: string | null;
};

type RequestHandler = (request: RequestType, meta: { service: string, method: string }) => Promise<unknown>;

export function createAuthServiceClient(handler: RequestHandler): AuthService {
    return {
      ExchangeToken(request) { // eslint-disable-line @typescript-eslint/no-unused-vars
        if (!request.tokenKey) {
          throw new Error("missing required field request.tokenKey");
        }
        const path = `auth/token/${request.tokenKey}`; // eslint-disable-line quotes
        const body = null;
        const queryParams: string[] = [];
        let uri = path;
        if (queryParams.length > 0) {
          uri += `?${queryParams.join("&")}`
        }   
        return handler({
          path: uri,
          method: "GET",
          body,
        }, {
          service: "AuthService",
          method: "ExchangeToken",
        }) as Promise<ExchangeTokenResponse>;
      },
      CheckToken(request) { // eslint-disable-line @typescript-eslint/no-unused-vars
        const path = `auth/check/token`; // eslint-disable-line quotes
        const body = null;
        const queryParams: string[] = [];
        let uri = path;
        if (queryParams.length > 0) {
          uri += `?${queryParams.join("&")}`
        }   
        return handler({
          path: uri,
          method: "GET",
          body,
        }, {
          service: "AuthService",
          method: "CheckToken",
        }) as Promise<CheckTokenResponse>;
      }
  }
}