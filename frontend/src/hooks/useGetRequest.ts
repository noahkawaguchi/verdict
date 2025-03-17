import { useState, useCallback } from 'react';
import { backendUrl } from '../config';

/**
 * Custom hook for making GET requests to the backend with a path parameter.
 *
 * @template TResponse - The expected response type.
 * @param endpoint - The API endpoint. Do not include the base URL, leading slash, or
 *                   parameters. Correct endpoint example: "poll"
 * @returns { data, error, loading, sendRequest }
 *          data - The expected TResponse or null.
 *          error - Any error encountered or null.
 *          loading - Whether the request is currently in progress.
 *          sendRequest - The function to trigger the request.
 */
const useGetRequest = <TResponse>(endpoint: string) => {
  const [data, setData] = useState<TResponse | null>(null);
  const [error, setError] = useState<Error | null>(null);
  const [loading, setLoading] = useState(false);

  /**
   * Sends a GET request using the specified path parameter.
   * 
   * @param pathParameter - The desired path parameter. Do not include a slash.
   */
  const sendRequest = useCallback(
    async (pathParameter: string) => {
      setLoading(true);
      setError(null);

      try {
        const response = await fetch(`${backendUrl}/${endpoint}/${pathParameter}`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        });
        if (!response.ok) {
          const errorText = await response.text();
          throw new Error(`HTTP ${response.status}: ${errorText}`);
        }
        const responseData = await response.json();
        setData(responseData);
      } catch (err) {
        setError(err instanceof Error ? err : new Error('Unknown error'));
      } finally {
        setLoading(false);
      }
    },
    [endpoint],
  );

  return { data, error, loading, sendRequest };
};

export default useGetRequest;
