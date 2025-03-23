import { renderHook } from '@testing-library/react';
import useMutationRequest from './useMutationRequest';
import { backendUrl } from '../config';
import { act } from 'react';

describe('useMutationRequest', () => {
  beforeEach(() => (globalThis.fetch = vi.fn()));

  afterEach(() => vi.resetAllMocks());

  afterAll(() => vi.restoreAllMocks());

  it('should set the initial return values correctly', () => {
    const { result } = renderHook(() =>
      useMutationRequest<Record<string, string>, Record<string, string>>('dummy-endpoint', 'POST'),
    );
    expect(result.current.data).toBeNull();
    expect(result.current.error).toBeNull();
    expect(result.current.loading).toEqual(false);
  });

  it('should call fetch with the correct arguments', async () => {
    const mockFetch = globalThis.fetch as ReturnType<typeof vi.fn>;
    const { result } = renderHook(() =>
      useMutationRequest<{ field1: string; field2: string }, Record<string, string>>(
        'dummy-endpoint',
        'POST',
      ),
    );
    await act(async () =>
      result.current.sendRequest({ field1: 'some data', field2: 'other data' }),
    );
    expect(mockFetch).toHaveBeenCalledExactlyOnceWith(`${backendUrl}/dummy-endpoint`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: `{"field1":"some data","field2":"other data"}`,
    });
  });

  it('should correctly handle a successful response', async () => {
    const mockFetch = globalThis.fetch as ReturnType<typeof vi.fn>;
    mockFetch.mockResolvedValueOnce({ ok: true, json: async () => ({ message: 'success!' }) });
    const { result } = renderHook(() =>
      useMutationRequest<{ field1: string; field2: string }, Record<string, string>>(
        'dummy-endpoint',
        'POST',
      ),
    );
    await act(async () =>
      result.current.sendRequest({ field1: 'some data', field2: 'other data' }),
    );
    expect(result.current.loading).toEqual(false);
    expect(result.current.error).toBeNull();
    expect(result.current.data).toEqual({ message: 'success!' });
  });

  it('should correctly handle an error response', async () => {
    const mockFetch = globalThis.fetch as ReturnType<typeof vi.fn>;
    mockFetch.mockResolvedValueOnce({
      ok: false,
      json: async () => ({ error: 'something went wrong' }),
    });
    const { result } = renderHook(() =>
      useMutationRequest<{ field1: string; field2: string }, Record<string, string>>(
        'dummy-endpoint',
        'POST',
      ),
    );
    await act(async () =>
      result.current.sendRequest({ field1: 'some data', field2: 'other data' }),
    );
    expect(result.current.loading).toEqual(false);
    expect(result.current.data).toBeNull();
    expect(result.current.error?.message).toEqual('something went wrong');
  });
});
