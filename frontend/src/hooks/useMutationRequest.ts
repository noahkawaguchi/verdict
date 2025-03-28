import { useState, useCallback, useRef } from 'react';
import { backendUrl } from '../config';

/**
 * Custom hook for making mutation requests to the backend (e.g. POST).
 *
 * @template TRequest - The type of the request body.
 * @template TResponse - The expected response type.
 * @param endpoint - The API endpoint. Do not include the base URL, leading slash, or
 *                   parameters. Correct endpoint example: "poll"
 * @returns { data, error, loading, sendRequest }
 *          data - The expected TResponse or null.
 *          error - Any error encountered or null.
 *          loading - Whether the request is currently in progress.
 *          sendRequest - The function to trigger the request.
 */
const useMutationRequest = <TRequest, TResponse>(
  endpoint: string,
  method: 'POST' | 'PUT' | 'PATCH',
) => {
  const [data, setData] = useState<TResponse | null>(null);
  const [error, setError] = useState<Error | null>(null);
  const [loading, setLoading] = useState(false);

  // Ref to persist across rerenders and avoid double requests, especially when using React's
  // StrictMode in development
  const isSubmitting = useRef(false);

  /**
   * Sends a mutation request using the specified body.
   *
   * @param body - The body to be sent in the request.
   */
  const sendRequest = useCallback(
    async (body: TRequest) => {
      if (isSubmitting.current) return;
      isSubmitting.current = true;

      setData(null);
      setError(null);
      setLoading(true);

      try {
        const response = await fetch(`${backendUrl}/${endpoint}`, {
          method: method,
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        });
        const parsedBody = await response.json();
        if (response.ok) setData(parsedBody);
        else throw new Error(parsedBody.error);
      } catch (err) {
        setError(err instanceof Error ? err : new Error('unknown error'));
      } finally {
        setLoading(false);
        isSubmitting.current = false;
      }
    },
    [endpoint, method],
  );

  return { data, error, loading, sendRequest };
};

export default useMutationRequest;
