import { useState, useCallback } from 'react';
import { backendUrl } from '../config';

const useGetRequest = <TResponse>(endpoint: string) => {
  const [data, setData] = useState<TResponse | null>(null);
  const [error, setError] = useState<Error | null>(null);
  const [loading, setLoading] = useState(false);

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
