import { useState, useCallback } from 'react';

const backendUrl = 'http://127.0.0.1:3000';

const useApiRequest = <TRequest, TResponse>(endpoint: string, method: 'POST' | 'PUT') => {
  const [data, setData] = useState<TResponse | null>(null);
  const [error, setError] = useState<Error | null>(null);
  const [loading, setLoading] = useState(false);

  const sendRequest = useCallback(
    async (body: TRequest) => {
      setLoading(true);
      setError(null);

      try {
        const response = await fetch(`${backendUrl}/${endpoint}`, {
          method: method,
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(body),
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
    [endpoint, method],
  );

  return { data, error, loading, sendRequest };
};

export default useApiRequest;
