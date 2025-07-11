//
// THIS CODE AND INFORMATION IS PROVIDED "AS IS" WITHOUT WARRANTY OF
// ANY KIND, EITHER EXPRESSED OR IMPLIED, INCLUDING BUT NOT LIMITED TO
// THE IMPLIED WARRANTIES OF MERCHANTABILITY AND/OR FITNESS FOR A
// PARTICULAR PURPOSE.
//
// Copyright (c) Microsoft Corporation. All rights reserved.
//
// Standard dll required functions and class factory implementation.

#include <windows.h>
#include <unknwn.h>
#include "Dll.h"
#include "helpers.h"
#include "CSampleProviderFilter.h" // Thêm include cho filter
#include "guid.h" // Đảm bảo có CLSID_CSampleProviderFilter

static long g_cRef = 0;   // global dll reference count
HINSTANCE g_hinst = NULL; // global dll hinstance

extern HRESULT CSample_CreateInstance(__in REFIID riid, __deref_out void** ppv);
extern HRESULT CSampleProviderFilter_CreateInstance(__in REFIID riid, __deref_out void** ppv);

class CClassFactory : public IClassFactory
{
public:
    CClassFactory(REFCLSID clsid) : _cRef(1), _clsid(clsid) {}

    // IUnknown
    IFACEMETHODIMP QueryInterface(__in REFIID riid, __deref_out void** ppv)
    {
        static const QITAB qit[] =
        {
            QITABENT(CClassFactory, IClassFactory),
            { 0 },
        };
        return QISearch(this, qit, riid, ppv);
    }

    IFACEMETHODIMP_(ULONG) AddRef()
    {
        return InterlockedIncrement(&_cRef);
    }

    IFACEMETHODIMP_(ULONG) Release()
    {
        long cRef = InterlockedDecrement(&_cRef);
        if (!cRef)
            delete this;
        return cRef;
    }

    // IClassFactory
    IFACEMETHODIMP CreateInstance(__in IUnknown* pUnkOuter, __in REFIID riid, __deref_out void** ppv)
    {
        if (pUnkOuter)
        {
            *ppv = NULL;
            return CLASS_E_NOAGGREGATION;
        }
        if (_clsid == CLSID_CSample)
        {
            return CSample_CreateInstance(riid, ppv);
        }
        else if (_clsid == CLSID_CSampleProviderFilter)
        {
            return CSampleProviderFilter_CreateInstance(riid, ppv);
        }
        *ppv = NULL;
        return CLASS_E_CLASSNOTAVAILABLE;
    }

    IFACEMETHODIMP LockServer(__in BOOL bLock)
    {
        if (bLock)
        {
            DllAddRef();
        }
        else
        {
            DllRelease();
        }
        return S_OK;
    }

private:
    ~CClassFactory() {}
    long _cRef;
    CLSID _clsid;
};

HRESULT CClassFactory_CreateInstance(__in REFCLSID rclsid, __in REFIID riid, __deref_out void** ppv)
{
    *ppv = NULL;
    if (rclsid == CLSID_CSample || rclsid == CLSID_CSampleProviderFilter)
    {
        CClassFactory* pcf = new CClassFactory(rclsid);
        if (pcf)
        {
            HRESULT hr = pcf->QueryInterface(riid, ppv);
            pcf->Release();
            return hr;
        }
        return E_OUTOFMEMORY;
    }
    return CLASS_E_CLASSNOTAVAILABLE;
}

void DllAddRef()
{
    InterlockedIncrement(&g_cRef);
}

void DllRelease()
{
    InterlockedDecrement(&g_cRef);
}

STDAPI DllCanUnloadNow()
{
    return (g_cRef > 0) ? S_FALSE : S_OK;
}

STDAPI DllGetClassObject(__in REFCLSID rclsid, __in REFIID riid, __deref_out void** ppv)
{
    return CClassFactory_CreateInstance(rclsid, riid, ppv);
}

STDAPI_(BOOL) DllMain(__in HINSTANCE hinstDll, __in DWORD dwReason, __in void*)
{
    switch (dwReason)
    {
    case DLL_PROCESS_ATTACH:
        DisableThreadLibraryCalls(hinstDll);
        break;
    case DLL_PROCESS_DETACH:
    case DLL_THREAD_ATTACH:
    case DLL_THREAD_DETACH:
        break;
    }

    g_hinst = hinstDll;
    return TRUE;
}

