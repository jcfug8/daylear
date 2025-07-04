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
    // the file to ocr
    file: File | undefined;
};

// the auth service
export interface FileService {
    // exchange a token key for a token
    UploadRecipeImage(request: UploadRecipeImageRequest): Promise<UploadRecipeImageResponse>;
    // ocr a recipe image
    OCRRecipe(request: OCRRecipeRequest): Promise<OCRRecipeResponse>;
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
      OCRRecipe(request) { // eslint-disable-line @typescript-eslint/no-unused-vars
        if (!request.file) {
          throw new Error("missing required field request.file");
        }
        const path = `files/meals/v1alpha1/recipes:ocr`; // eslint-disable-line quotes
        const body = request.file;
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