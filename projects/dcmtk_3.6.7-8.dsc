-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA512

Format: 3.0 (quilt)
Source: dcmtk
Binary: dcmtk, libdcmtk17, libdcmtk-dev, dcmtk-doc
Architecture: any all
Version: 3.6.7-8
Maintainer: Debian Med Packaging Team <debian-med-packaging@lists.alioth.debian.org>
Uploaders: Gert Wollny <gewo@debian.org>, Mathieu Malaterre <malat@debian.org>
Homepage: https://dicom.offis.de/dcmtk
Standards-Version: 4.6.2
Vcs-Browser: https://salsa.debian.org/med-team/dcmtk
Vcs-Git: https://salsa.debian.org/med-team/dcmtk.git
Testsuite: autopkgtest
Testsuite-Triggers: python3-pydicom
Build-Depends: cmake, debhelper-compat (= 13), gettext, help2man, libpng-dev, libsndfile1-dev, libssl-dev, libtiff-dev, libwrap0-dev, libxml2-dev, zlib1g-dev
Build-Depends-Indep: doxygen, graphviz
Package-List:
 dcmtk deb science optional arch=any
 dcmtk-doc deb doc optional arch=all
 libdcmtk-dev deb libdevel optional arch=any
 libdcmtk17 deb libs optional arch=any
Checksums-Sha1:
 d6304e564458d97074c48bd8054c30d9a8bf1062 6641083 dcmtk_3.6.7.orig.tar.gz
 54978fa757c4a6bb007f8f4569c0f6d7518ccfe4 40772 dcmtk_3.6.7-8.debian.tar.xz
Checksums-Sha256:
 17705dcdb2047d1266bb4e92dbf4aa6d4967819e8e3e94f39b7df697661b4860 6641083 dcmtk_3.6.7.orig.tar.gz
 f7c8ebae7050cf99391688cf321f64a1d3794c9da36eeae687e2b988cd341447 40772 dcmtk_3.6.7-8.debian.tar.xz
Files:
 ef8323ad0d9067a035af238435d1948d 6641083 dcmtk_3.6.7.orig.tar.gz
 bcd22e8074fe1ef11add7543cdbff067 40772 dcmtk_3.6.7-8.debian.tar.xz

-----BEGIN PGP SIGNATURE-----

iQJFBAEBCgAvFiEEaTNn/67NjqrNHwY7AXHhgorgk0UFAmO1Ug8RHG1hbGF0QGRl
Ymlhbi5vcmcACgkQAXHhgorgk0Vu9w/+NuyLQ1A9JUo21HsaziSIH0Cq3X2omWNn
XpxJZbWMQ/CEVmu6FT4thidp1hI2uvnp/7CSei5/mClfv9fzamXI4yn3a5VAlpqh
c8IXbnHoxZNa/gP4IBrTA1/Af4rifga4b1JGl4rvT8cek7yp8JFkaHQw6/sdUgqN
TcSNTihSfUDI2p3pNnyi/AuATSFm0FT6In1z8IjpiITs4cjtg9WlSRs7JJi7LLSc
nktqPL7s6Tirv6wzOm9072lBd3ahiD7C7e3QktSrgAjrZTcDmun3GxxHB+KyEFKi
Ljt2+eAcUgHSTqiSLyFMW3ZoenGQABGF6FCYoXTobvlofthIulq15GAtR6a43na6
oczSfqM1z+PUAhxGyWc0QyeIaI1GqmYHoJJEuYBuoVcZiG9OWTcU9G8DCQdiOcIp
5Gt+bIKrBsXnwQbeGtaZKMv8bDN6FGySUcigmsxpdW+ui5fqwL/5UQ1/d3Mq/824
6XKposL1q+gOxHDF9jfCi5LhhWMIJL9LPEKAnzdEsVlrRIes0v2mh7J8pHkEo3Z/
4NF72C1R4ymyf5T7HoYfF6PT3lO4M2neJ4VH3nb3F2VlEX/OBvhLft0NU29USSpf
JutJbmX3lgWi6o9U8WCaK1jyHYwb4B0jIM8J/b0oEWztTvaKqmGDfeyk68Trh3DS
JRXITYmk47k=
=oWRk
-----END PGP SIGNATURE-----
