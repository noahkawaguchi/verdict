import { useState, useCallback, useRef } from 'react';
import { backendUrl } from '../config';

const useMutationRequest = <TRequest, TResponse>(endpoint: string, method: 'POST' | 'PUT') => {
  const [data, setData] = useState<TResponse | null>(null);
  const [error, setError] = useState<Error | null>(null);
  const [loading, setLoading] = useState(false);

  // Ref to persist across rerenders and avoid double requests, especially when using React's
  // StrictMode in development
  const isSubmitting = useRef(false);

  const sendRequest = useCallback(
    async (body: TRequest) => {
      if (isSubmitting.current) return;
      isSubmitting.current = true;

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
        isSubmitting.current = false;
      }
    },
    [endpoint, method],
  );

  return { data, error, loading, sendRequest };
};

export default useMutationRequest;
