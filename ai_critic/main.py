import json

def verdict_for(game):
    price = game.get("price", 0)
    discount = game.get("discount", 0)

    if price == 0:
        return "Игра бесплатная! Забирай"
    if discount >= 70:
        return "Огромная скидка! Надо брать"
    if discount > 0:
        return "Хорошая скидка, можно подумать"
    return "Жди распродажи, скидки нет"

with open("../games_data.json", encoding="utf-8") as f:
    games = json.load(f)

if games is None:
    games = []

analyzed = []
for game in games:
    if game is None:
        continue
    game["verdict"] = verdict_for(game)
    analyzed.append(game)

with open("../analyzed_games.json", "w", encoding="utf-8") as f:
    json.dump(analyzed, f, ensure_ascii=False, indent=2)

print(f"Analyzed {len(analyzed)} games")
