import { renderHook } from '@testing-library/react';
import useGetRequest from './useGetRequest';
import { backendUrl } from '../config';
import { act } from 'react';

describe('useGetRequest', () => {
  beforeEach(() => (globalThis.fetch = vi.fn()));

  afterEach(() => vi.resetAllMocks());

  afterAll(() => vi.restoreAllMocks());

  it('should set the initial return values correctly', () => {
    const { result } = renderHook(() => useGetRequest<Record<string, string>>('dummy-endpoint'));
    expect(result.current.data).toBeNull();
    expect(result.current.error).toBeNull();
    expect(result.current.loading).toEqual(false);
  });

  it('should call fetch with the correct arguments', async () => {
    const mockFetch = globalThis.fetch as ReturnType<typeof vi.fn>;
    const { result } = renderHook(() => useGetRequest<Record<string, string>>('dummy-endpoint'));
    await act(async () => result.current.sendRequest('dummy-parameter'));
    expect(mockFetch).toHaveBeenCalledExactlyOnceWith(
      `${backendUrl}/dummy-endpoint/dummy-parameter`,
      { method: 'GET', headers: { 'Content-Type': 'application/json' } },
    );
  });

  it('should correctly handle a successful response', async () => {
    const mockFetch = globalThis.fetch as ReturnType<typeof vi.fn>;
    mockFetch.mockResolvedValueOnce({ ok: true, json: async () => ({ message: 'success!' }) });
    const { result } = renderHook(() => useGetRequest<{ message: string }>('dummy-endpoint'));
    await act(async () => result.current.sendRequest('dummy-parameter'));
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
    const { result } = renderHook(() => useGetRequest<{ message: string }>('dummy-endpoint'));
    await act(async () => result.current.sendRequest('dummy-parameter'));
    expect(result.current.loading).toEqual(false);
    expect(result.current.data).toBeNull();
    expect(result.current.error?.message).toEqual('something went wrong');
  });
});
