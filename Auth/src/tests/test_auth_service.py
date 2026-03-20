"""Unit-тесты для services/auth.py — не требуют БД."""
from services.auth import get_password_hash, verify_password


def test_hash_returns_string():
    assert isinstance(get_password_hash("password"), str)


def test_hash_differs_from_plain():
    pw = "secret"
    assert get_password_hash(pw) != pw


def test_same_password_produces_different_hashes():
    # Argon2 использует случайную соль
    pw = "same"
    assert get_password_hash(pw) != get_password_hash(pw)


def test_verify_correct_password():
    pw = "correct_password"
    assert verify_password(pw, get_password_hash(pw)) is True


def test_verify_wrong_password():
    hashed = get_password_hash("correct")
    assert verify_password("wrong", hashed) is False
