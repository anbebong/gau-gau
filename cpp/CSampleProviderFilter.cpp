#include "CSampleProviderFilter.h"
#include "guid.h"
#include <new>

CSampleProviderFilter::CSampleProviderFilter() : _cRef(1) {}
CSampleProviderFilter::~CSampleProviderFilter() {}

IFACEMETHODIMP CSampleProviderFilter::QueryInterface(REFIID riid, void** ppv)
{
    if (!ppv) return E_POINTER;
    if (riid == IID_IUnknown || riid == IID_ICredentialProviderFilter)
    {
        *ppv = static_cast<ICredentialProviderFilter*>(this);
        AddRef();
        return S_OK;
    }
    *ppv = nullptr;
    return E_NOINTERFACE;
}

IFACEMETHODIMP_(ULONG) CSampleProviderFilter::AddRef()
{
    return InterlockedIncrement(&_cRef);
}

IFACEMETHODIMP_(ULONG) CSampleProviderFilter::Release()
{
    LONG cRef = InterlockedDecrement(&_cRef);
    if (!cRef)
        delete this;
    return cRef;
}

IFACEMETHODIMP CSampleProviderFilter::Filter(
    CREDENTIAL_PROVIDER_USAGE_SCENARIO /*cpus*/, DWORD /*dwFlags*/,
    GUID* rgclsidProviders, BOOL* rgbAllow, DWORD cProviders)
{
    for (DWORD i = 0; i < cProviders; i++)
    {
        if (IsEqualGUID(rgclsidProviders[i], CLSID_CSample))
            rgbAllow[i] = TRUE;
        else
            rgbAllow[i] = FALSE;
    }
    return S_OK;
}

IFACEMETHODIMP CSampleProviderFilter::UpdateRemoteCredential(
    const CREDENTIAL_PROVIDER_CREDENTIAL_SERIALIZATION* /*pcpcsIn*/,
    CREDENTIAL_PROVIDER_CREDENTIAL_SERIALIZATION* /*pcpcsOut*/)
{
    return E_NOTIMPL;
}

HRESULT CSampleProviderFilter_CreateInstance(REFIID riid, void** ppv)
{
    if (!ppv) return E_POINTER;
    *ppv = nullptr;
    CSampleProviderFilter* pObj = new (std::nothrow) CSampleProviderFilter();
    if (!pObj) return E_OUTOFMEMORY;
    HRESULT hr = pObj->QueryInterface(riid, ppv);
    pObj->Release();
    return hr;
}
