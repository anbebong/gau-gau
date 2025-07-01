#pragma once
#include <windows.h>
#include <credentialprovider.h>

class CSampleProviderFilter : public ICredentialProviderFilter
{
public:
    // IUnknown
    IFACEMETHODIMP QueryInterface(REFIID riid, void** ppv);
    IFACEMETHODIMP_(ULONG) AddRef();
    IFACEMETHODIMP_(ULONG) Release();

    // ICredentialProviderFilter
    IFACEMETHODIMP Filter(CREDENTIAL_PROVIDER_USAGE_SCENARIO cpus, DWORD dwFlags, GUID* rgclsidProviders, BOOL* rgbAllow, DWORD cProviders);
    IFACEMETHODIMP UpdateRemoteCredential(const CREDENTIAL_PROVIDER_CREDENTIAL_SERIALIZATION* pcpcsIn, CREDENTIAL_PROVIDER_CREDENTIAL_SERIALIZATION* pcpcsOut);

    CSampleProviderFilter();
    ~CSampleProviderFilter();

private:
    LONG _cRef;
};

HRESULT CSampleProviderFilter_CreateInstance(REFIID riid, void** ppv);
