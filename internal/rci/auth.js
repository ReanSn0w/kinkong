authenticate(e, i)
{
    return this.sendAuthRequest().pipe((0, ji.W)(o => (0, S.of)(o)), (0, se.n)(o => {
        if (200 === o.status)
            return this.isAuthenticated = !0, (0, S.of)(!0);
        const a = !o.headers.get("www-authenticate"),
            r=o?.headers?.get("x-ndm-isplock"),
            c = this.getDocumentCookies(tv);
        this.cookieStorageProvider.setItem(tv, c),
        this.isAuthAvailable$.next(!r || !a);
        const d = this.getEncryptedPassword({
            token: o.headers.get("X-NDM-Challenge") || "",
            realm: o.headers.get("X-NDM-Realm") || "",
            login: e,
            password: i
        });
        return this.httpClient.post(d_, {
            login: e,
            password: d
        }, {
            observe: "response"
        }).pipe((0, se.n)(() => this.checkIfAuthenticated()))
    }))
}

getEncryptedPassword(e)
{
    const {token: i, login: o, realm: a, password: r} = e,
        c = new rV.A("SHA-256", "TEXT"),
        d = ly.V.hashStr(`${o}:${a}:${r}`);
    return c.update(i + String(d)), c.getHash("HEX")
}

update(e)
{
    return this.writeUserSettings(e).pipe((0, g.T)(() => (this.userSettings$.next(e), e)))
}

writeUserSettings(e)
{
    const i = JSON.stringify(e);
    return this.systemEnvironmentSetApi.perform({
        name: this.ENV_KEY,
        value: i
    })
}