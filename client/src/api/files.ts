// the response to exchange a token key for a token
export type UploadRecipeImageResponse = {
    // the recipes
    imageUri: string | undefined;
};

// the request to exchange a token key for a token
export type UploadRecipeImageRequest = {
    // the recipe to update
    name: string | undefined;
    // the file to upload
    file: File | undefined;
};

// the auth service
export interface FileService {
    // exchange a token key for a token
    UploadRecipeImage(request: UploadRecipeImageRequest): Promise<UploadRecipeImageResponse>;
  }
  
type RequestType = {
    path: string;
    method: string;
    body: File | null;
};

type RequestHandler = (request: RequestType, meta: { service: string, method: string }) => Promise<unknown>;

export function createFileServiceClient(handler: RequestHandler): FileService {
    return {
      UploadRecipeImage(request) { // eslint-disable-line @typescript-eslint/no-unused-vars
        if (!request.name) {
          throw new Error("missing required field request.name");
        }
        if (!request.file) {
          throw new Error("missing required field request.file");
        }
        const path = `files/meals/v1alpha1/${request.name}/image`; // eslint-disable-line quotes
        const body = request.file;
        const queryParams: string[] = [];
        let uri = path;
        if (queryParams.length > 0) {
          uri += `?${queryParams.join("&")}`
        }   
        return handler({
          path: uri,
          method: "PUT",
          body,
        }, {
          service: "FileService",
          method: "UploadRecipeImage",
        }) as Promise<UploadRecipeImageResponse>;
      },
  }
}