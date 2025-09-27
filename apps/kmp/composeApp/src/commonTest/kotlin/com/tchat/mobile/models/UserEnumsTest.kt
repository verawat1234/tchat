package com.tchat.mobile.models

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNull
import kotlin.test.assertTrue

class UserEnumsTest {

    @Test
    fun testUserStatusFromValue() {
        assertEquals(UserStatus.ACTIVE, UserStatus.fromValue("active"))
        assertEquals(UserStatus.SUSPENDED, UserStatus.fromValue("suspended"))
        assertEquals(UserStatus.DELETED, UserStatus.fromValue("deleted"))
        assertNull(UserStatus.fromValue("invalid"))
    }

    @Test
    fun testUserStatusValues() {
        assertEquals("active", UserStatus.ACTIVE.value)
        assertEquals("suspended", UserStatus.SUSPENDED.value)
        assertEquals("deleted", UserStatus.DELETED.value)
    }

    @Test
    fun testKYCTierFromValue() {
        assertEquals(KYCTier.UNVERIFIED, KYCTier.fromValue(0))
        assertEquals(KYCTier.BASIC, KYCTier.fromValue(1))
        assertEquals(KYCTier.STANDARD, KYCTier.fromValue(2))
        assertEquals(KYCTier.PREMIUM, KYCTier.fromValue(3))
        assertNull(KYCTier.fromValue(99))
    }

    @Test
    fun testKYCTierValues() {
        assertEquals(0, KYCTier.UNVERIFIED.value)
        assertEquals(1, KYCTier.BASIC.value)
        assertEquals(2, KYCTier.STANDARD.value)
        assertEquals(3, KYCTier.PREMIUM.value)
    }

    @Test
    fun testKYCTierDisplayNames() {
        assertEquals("Unverified", KYCTier.UNVERIFIED.displayName)
        assertEquals("Basic", KYCTier.BASIC.displayName)
        assertEquals("Standard", KYCTier.STANDARD.displayName)
        assertEquals("Premium", KYCTier.PREMIUM.displayName)
    }

    @Test
    fun testCountryFromCode() {
        assertEquals(Country.THAILAND, Country.fromCode("TH"))
        assertEquals(Country.SINGAPORE, Country.fromCode("SG"))
        assertEquals(Country.INDONESIA, Country.fromCode("ID"))
        assertEquals(Country.MALAYSIA, Country.fromCode("MY"))
        assertEquals(Country.PHILIPPINES, Country.fromCode("PH"))
        assertEquals(Country.VIETNAM, Country.fromCode("VN"))
        assertNull(Country.fromCode("XX"))
    }

    @Test
    fun testCountryValues() {
        assertEquals("TH", Country.THAILAND.code)
        assertEquals("Thailand", Country.THAILAND.displayName)

        assertEquals("SG", Country.SINGAPORE.code)
        assertEquals("Singapore", Country.SINGAPORE.displayName)

        assertEquals("ID", Country.INDONESIA.code)
        assertEquals("Indonesia", Country.INDONESIA.displayName)

        assertEquals("MY", Country.MALAYSIA.code)
        assertEquals("Malaysia", Country.MALAYSIA.displayName)

        assertEquals("PH", Country.PHILIPPINES.code)
        assertEquals("Philippines", Country.PHILIPPINES.displayName)

        assertEquals("VN", Country.VIETNAM.code)
        assertEquals("Vietnam", Country.VIETNAM.displayName)
    }

    @Test
    fun testCountryGetAllCountries() {
        val allCountries = Country.getAllCountries()
        assertEquals(6, allCountries.size)
        assertTrue(allCountries.contains(Country.THAILAND))
        assertTrue(allCountries.contains(Country.SINGAPORE))
        assertTrue(allCountries.contains(Country.INDONESIA))
        assertTrue(allCountries.contains(Country.MALAYSIA))
        assertTrue(allCountries.contains(Country.PHILIPPINES))
        assertTrue(allCountries.contains(Country.VIETNAM))
    }

    @Test
    fun testVerificationTierFromValue() {
        assertEquals(VerificationTier.NONE, VerificationTier.fromValue(0))
        assertEquals(VerificationTier.PHONE, VerificationTier.fromValue(1))
        assertEquals(VerificationTier.EMAIL, VerificationTier.fromValue(2))
        assertEquals(VerificationTier.KYC, VerificationTier.fromValue(3))
        assertEquals(VerificationTier.FULL, VerificationTier.fromValue(4))
        assertNull(VerificationTier.fromValue(99))
    }

    @Test
    fun testVerificationTierValues() {
        assertEquals(0, VerificationTier.NONE.value)
        assertEquals(1, VerificationTier.PHONE.value)
        assertEquals(2, VerificationTier.EMAIL.value)
        assertEquals(3, VerificationTier.KYC.value)
        assertEquals(4, VerificationTier.FULL.value)
    }

    @Test
    fun testVerificationTierDisplayNames() {
        assertEquals("No Verification", VerificationTier.NONE.displayName)
        assertEquals("Phone Verified", VerificationTier.PHONE.displayName)
        assertEquals("Email Verified", VerificationTier.EMAIL.displayName)
        assertEquals("KYC Verified", VerificationTier.KYC.displayName)
        assertEquals("Full Verification", VerificationTier.FULL.displayName)
    }
}