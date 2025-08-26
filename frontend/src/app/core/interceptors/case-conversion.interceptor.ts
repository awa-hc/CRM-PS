import { HttpInterceptorFn, HttpRequest, HttpResponse } from '@angular/common/http';
import { map } from 'rxjs/operators';

export const caseConversionInterceptor: HttpInterceptorFn = (req, next) => {
  // Convert request body from camelCase to snake_case
  let modifiedReq = req;
  
  if (req.body && typeof req.body === 'object' && !isFormData(req.body)) {
    const convertedBody = convertKeysToSnakeCase(req.body);
    modifiedReq = req.clone({ body: convertedBody });
  }

  return next(modifiedReq).pipe(
    map(event => {
      // Convert response body from snake_case to camelCase
      if (event instanceof HttpResponse && event.body && typeof event.body === 'object') {
        const convertedBody = convertKeysToCamelCase(event.body);
        return event.clone({ body: convertedBody });
      }
      return event;
    })
  );
};

function isFormData(obj: any): boolean {
  return obj instanceof FormData;
}

function convertKeysToSnakeCase(obj: any): any {
  if (obj === null || typeof obj !== 'object') {
    return obj;
  }

  if (Array.isArray(obj)) {
    return obj.map(convertKeysToSnakeCase);
  }

  const converted: any = {};
  for (const key in obj) {
    if (obj.hasOwnProperty(key)) {
      const snakeKey = camelToSnakeCase(key);
      converted[snakeKey] = convertKeysToSnakeCase(obj[key]);
    }
  }
  return converted;
}

function convertKeysToCamelCase(obj: any): any {
  if (obj === null || typeof obj !== 'object') {
    return obj;
  }

  if (Array.isArray(obj)) {
    return obj.map(convertKeysToCamelCase);
  }

  const converted: any = {};
  for (const key in obj) {
    if (obj.hasOwnProperty(key)) {
      const camelKey = snakeToCamelCase(key);
      converted[camelKey] = convertKeysToCamelCase(obj[key]);
    }
  }
  return converted;
}

function camelToSnakeCase(str: string): string {
  return str.replace(/[A-Z]/g, letter => `_${letter.toLowerCase()}`);
}

function snakeToCamelCase(str: string): string {
  return str.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase());
}