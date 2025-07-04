import type { Recipe } from "@/genapi/api/meals/recipe/v1alpha1";

// the response to upload a recipe image
export type UploadRecipeImageResponse = {
    // the recipes
    imageUri: string | undefined;
};

// the request to upload a recipe image
export type UploadRecipeImageRequest = {
    // the recipe to update
    name: string | undefined;
    // the file to upload
    file: File | undefined;
};

// the response to ocr a recipe image
export type OCRRecipeResponse = {
    // the recipe
    recipe: Recipe | undefined;
};

// the request to ocr a recipe image
export type OCRRecipeRequest = {
    // the files to ocr
    files: File[] | undefined;
};

// the response to upload a circle image
export type UploadCircleImageResponse = {
    imageUri: string | undefined;
};

// the request to upload a circle image
export type UploadCircleImageRequest = {
    name: string | undefined;
    file: File | undefined;
};

// the auth service
export interface FileService {
    // upload a recipe image
    UploadRecipeImage(request: UploadRecipeImageRequest): Promise<UploadRecipeImageResponse>;
    // upload a circle image
    UploadCircleImage(request: UploadCircleImageRequest): Promise<UploadCircleImageResponse>;
    // ocr a recipe image
    OCRRecipe(request: OCRRecipeRequest): Promise<OCRRecipeResponse>;
  }
  
type RequestType = {
    path: string;
    method: string;
    body: File | FormData | null;
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
        const body = new FormData();
        body.append("file", request.file);
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
      UploadCircleImage(request) {
        if (!request.name) {
          throw new Error("missing required field request.name");
        }
        if (!request.file) {
          throw new Error("missing required field request.file");
        }
        const path = `files/circles/v1alpha1/${request.name}/image`;
        const body = new FormData();
        body.append("file", request.file);
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
          method: "UploadCircleImage",
        }) as Promise<UploadCircleImageResponse>;
      },
      OCRRecipe(request) { // eslint-disable-line @typescript-eslint/no-unused-vars
        if (!request.files) {
          throw new Error("missing required field request.files");
        }
        const path = `files/meals/v1alpha1/recipes:ocr`; // eslint-disable-line quotes
        const body = new FormData();
        for (const file of request.files) {
            body.append("files", file);
        }
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
          service: "FileService",
          method: "OCRRecipe",
        }) as Promise<OCRRecipeResponse>;
      }
  }
}