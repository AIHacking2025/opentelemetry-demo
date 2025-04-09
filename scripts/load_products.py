#!/usr/bin/env python3

import json
import os
from dataclasses import dataclass
from pathlib import Path

import psycopg
from dotenv import load_dotenv

# Get the project root directory (one level up from this script)
PROJECT_ROOT = Path(__file__).parent.parent


@dataclass
class PriceUsd:
    currency_code: str
    units: int
    nanos: int


@dataclass
class Product:
    id: str
    name: str
    description: str
    picture: str
    price_usd: PriceUsd
    categories: list[str]

    @classmethod
    def from_dict(cls, data: dict) -> "Product":
        return cls(
            id=data["id"],
            name=data["name"],
            description=data["description"],
            picture=data["picture"],
            price_usd=PriceUsd(**data["priceUsd"]),
            categories=data["categories"],
        )


def load_products() -> list[Product]:
    products_file = (
        PROJECT_ROOT / "src" / "product-catalog" / "products" / "products.json"
    )
    with open(products_file) as f:
        data = json.load(f)
        return [Product.from_dict(p) for p in data["products"]]


def main() -> None:
    # Load environment variables from .env.secret in project root
    env_file = PROJECT_ROOT / ".env.secret"
    if not env_file.exists():
        raise FileNotFoundError(f"Environment file not found at {env_file}")

    load_dotenv(env_file)
    db_url = os.environ["NEON_DB_URL"]

    # Connect to the database
    with psycopg.connect(db_url) as conn:
        with conn.cursor() as cur:
            # First, clear existing products
            cur.execute("DELETE FROM products")

            # Insert products
            for product in load_products():
                cur.execute(
                    """
                    INSERT INTO products (
                        id, name, description, picture,
                        price_units, price_nanos, price_currency_code,
                        categories
                    )
                    VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
                    """,
                    (
                        product.id,
                        product.name,
                        product.description,
                        product.picture,
                        product.price_usd.units,
                        product.price_usd.nanos,
                        product.price_usd.currency_code,
                        product.categories,
                    ),
                )

        conn.commit()
        print("Successfully loaded product data into the database")


if __name__ == "__main__":
    main()
