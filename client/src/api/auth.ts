// the response to exchange a token key for a token
export type ExchangeTokenResponse = {
    // the recipes
    token: String | undefined;
};

// the request to exchange a token key for a token
export type ExchangeTokenRequest = {
    // the recipe to update
    tokenKey: String | undefined;
    // the fields to update
};

// the recipe service
export interface AuthService {
    // create a recipe
    ExchangeToken(request: ExchangeTokenRequest): Promise<ExchangeTokenResponse>;
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
          method: "POST",
          body,
        }, {
          service: "RecipeService",
          method: "ExchangeToken",
        }) as Promise<ExchangeTokenResponse>;
      }
  }
}