<?php
// recipe_catalog.php — PHP версия

$dataFile = 'recipes.json';

function loadRecipes() {
    global $dataFile;
    if (file_exists($dataFile)) {
        $json = file_get_contents($dataFile);
        return json_decode($json, true) ?: [];
    }
    return [];
}

function saveRecipes($recipes) {
    global $dataFile;
    file_put_contents($dataFile, json_encode($recipes, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE));
}

$recipes = loadRecipes();

function color($text, $code) {
    return "\033[{$code}m{$text}\033[0m";
}

function listAll($recipes) {
    if (empty($recipes)) {
        echo color("Каталог пуст.\n", '33');
        return;
    }
    printf(color("%-4s %-25s %-12s %-8s %-10s %-8s\n", '36'), "ID", "Название", "Категория", "Время", "Сложность", "Рейтинг");
    echo str_repeat("-", 75) . "\n";
    foreach ($recipes as $r) {
        $ratingStr = $r['ratings_count'] > 0 ? round($r['rating'], 1) . "★" : "—";
        $diffColor = $r['difficulty'] == 'лёгкая' ? '32' : ($r['difficulty'] == 'средняя' ? '33' : '31');
        printf("%-4d %-25s %-12s %-8d %s%-10s\033[0m %-8s\n", $r['id'], $r['name'], $r['category'], $r['prep_time'], color("", $diffColor), $r['difficulty'], $ratingStr);
    }
}

function searchRecipes($recipes, $query) {
    $query = strtolower($query);
    $results = array_filter($recipes, function($r) use ($query) {
        if (stripos($r['name'], $query) !== false) return true;
        foreach ($r['ingredients'] as $ing) {
            if (stripos($ing, $query) !== false) return true;
        }
        return false;
    });
    if (empty($results)) {
        echo color("Ничего не найдено.\n", '33');
        return;
    }
    foreach ($results as $r) {
        echo "{$r['id']}: {$r['name']} | {$r['category']} | {$r['prep_time']} мин | {$r['difficulty']} | " . round($r['rating'], 1) . "★\n";
    }
}

function filterRecipes($recipes, $field, $value) {
    $results = [];
    foreach ($recipes as $r) {
        if ($field == 'category' && strcasecmp($r['category'], $value) == 0) {
            $results[] = $r;
        } elseif ($field == 'difficulty' && strcasecmp($r['difficulty'], $value) == 0) {
            $results[] = $r;
        } elseif ($field == 'prep_time' && $r['prep_time'] <= (int)$value) {
            $results[] = $r;
        }
    }
    if (empty($results)) {
        echo color("Нет рецептов, соответствующих фильтру.\n", '33');
        return;
    }
    foreach ($results as $r) {
        echo "{$r['id']}: {$r['name']} | {$r['category']} | {$r['prep_time']} мин | {$r['difficulty']}\n";
    }
}

function sortRecipes(&$recipes, $field, $reverse) {
    usort($recipes, function($a, $b) use ($field, $reverse) {
        $cmp = 0;
        if ($field == 'name') $cmp = strcasecmp($a['name'], $b['name']);
        elseif ($field == 'prep_time') $cmp = $a['prep_time'] - $b['prep_time'];
        elseif ($field == 'rating') $cmp = $a['rating'] <=> $b['rating'];
        else return 0;
        return $reverse ? -$cmp : $cmp;
    });
    listAll($recipes);
}

function deleteRecipe(&$recipes, $id) {
    foreach ($recipes as $i => $r) {
        if ($r['id'] == $id) {
            array_splice($recipes, $i, 1);
            saveRecipes($recipes);
            return true;
        }
    }
    return false;
}

function editRecipe(&$recipes, $id, $field, $value) {
    foreach ($recipes as &$r) {
        if ($r['id'] == $id) {
            switch ($field) {
                case 'name': $r['name'] = $value; break;
                case 'category': $r['category'] = $value; break;
                case 'ingredients': $r['ingredients'] = array_map('trim', explode(',', $value)); break;
                case 'instructions': $r['instructions'] = $value; break;
                case 'prep_time': $r['prep_time'] = (int)$value; break;
                case 'difficulty': $r['difficulty'] = $value; break;
                default: return false;
            }
            saveRecipes($recipes);
            return true;
        }
    }
    return false;
}

function rateRecipe(&$recipes, $id, $rating) {
    foreach ($recipes as &$r) {
        if ($r['id'] == $id) {
            $total = $r['rating'] * $r['ratings_count'] + $rating;
            $r['ratings_count']++;
            $r['rating'] = $total / $r['ratings_count'];
            saveRecipes($recipes);
            return true;
        }
    }
    return false;
}

function stats($recipes) {
    if (empty($recipes)) {
        echo "Нет данных.\n";
        return;
    }
    $categories = [];
    $difficulties = [];
    $totalRating = 0;
    $ratedCount = 0;
    foreach ($recipes as $r) {
        $categories[$r['category']] = ($categories[$r['category']] ?? 0) + 1;
        $difficulties[$r['difficulty']] = ($difficulties[$r['difficulty']] ?? 0) + 1;
        if ($r['ratings_count'] > 0) {
            $totalRating += $r['rating'];
            $ratedCount++;
        }
    }
    $avgRating = $ratedCount > 0 ? $totalRating / $ratedCount : 0;
    echo color("📊 Статистика:\n", '36');
    echo "  Всего рецептов: " . count($recipes) . "\n";
    echo "  Средний рейтинг: " . round($avgRating, 1) . "★ (из $ratedCount оценённых)\n";
    echo "  По категориям:\n";
    arsort($categories);
    foreach ($categories as $c => $count) {
        echo "    $c: $count\n";
    }
    echo "  По сложности:\n";
    foreach ($difficulties as $d => $count) {
        echo "    $d: $count\n";
    }
}

function exportJSON($recipes, $filename = 'recipes_export.json') {
    file_put_contents($filename, json_encode($recipes, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE));
    echo color("💾 Экспорт JSON: $filename\n", '32');
}

function exportCSV($recipes, $filename = 'recipes_export.csv') {
    if (empty($recipes)) {
        echo color("Нет данных для экспорта.\n", '33');
        return;
    }
    $fp = fopen($filename, 'w');
    fputcsv($fp, ['ID', 'Название', 'Категория', 'Ингредиенты', 'Инструкция', 'Время', 'Сложность', 'Рейтинг']);
    foreach ($recipes as $r) {
        fputcsv($fp, [$r['id'], $r['name'], $r['category'], implode(', ', $r['ingredients']), $r['instructions'], $r['prep_time'], $r['difficulty'], round($r['rating'], 1)]);
    }
    fclose($fp);
    echo color("💾 Экспорт CSV: $filename\n", '32');
}

function main() {
    global $recipes;
    while (true) {
        echo "\n" . color("🍰 Dessert Recipe Catalog (PHP)\n", '36');
        echo "1. Добавить рецепт\n";
        echo "2. Показать все рецепты\n";
        echo "3. Поиск рецептов\n";
        echo "4. Фильтрация\n";
        echo "5. Сортировка\n";
        echo "6. Удалить рецепт\n";
        echo "7. Редактировать рецепт\n";
        echo "8. Оценить рецепт\n";
        echo "9. Статистика\n";
        echo "10. Экспорт\n";
        echo "11. Выход\n";
        echo "Выберите действие: ";
        $choice = trim(fgets(STDIN));

        switch ($choice) {
            case '1':
                echo "Название: ";
                $name = trim(fgets(STDIN));
                echo "Категория (торт, пирожное, печенье, мороженое, другое): ";
                $category = trim(fgets(STDIN));
                echo "Ингредиенты (через запятую): ";
                $ingredients = array_map('trim', explode(',', trim(fgets(STDIN))));
                echo "Инструкция: ";
                $instructions = trim(fgets(STDIN));
                echo "Время приготовления (мин): ";
                $prepTime = (int) trim(fgets(STDIN));
                echo "Сложность (лёгкая/средняя/сложная): ";
                $difficulty = strtolower(trim(fgets(STDIN)));
                if (!in_array($difficulty, ['лёгкая', 'средняя', 'сложная'])) $difficulty = 'средняя';
                $id = count($recipes) + 1;
                $recipes[] = ['id' => $id, 'name' => $name, 'category' => $category, 'ingredients' => $ingredients,
                              'instructions' => $instructions, 'prep_time' => $prepTime, 'difficulty' => $difficulty,
                              'rating' => 0, 'ratings_count' => 0];
                saveRecipes($recipes);
                echo color("✅ Рецепт добавлен (ID: $id)\n", '32');
                break;
            case '2': listAll($recipes); break;
            case '3':
                echo "Введите запрос (название или ингредиент): ";
                $query = trim(fgets(STDIN));
                searchRecipes($recipes, $query);
                break;
            case '4':
                echo "Фильтровать по: category, difficulty, prep_time\n";
                echo "Поле: ";
                $field = trim(fgets(STDIN));
                echo "Значение: ";
                $value = trim(fgets(STDIN));
                filterRecipes($recipes, $field, $value);
                break;
            case '5':
                echo "Сортировать по: name, prep_time, rating\n";
                echo "Поле: ";
                $field = trim(fgets(STDIN));
                echo "По убыванию? (y/n): ";
                $reverse = trim(fgets(STDIN)) == 'y';
                sortRecipes($recipes, $field, $reverse);
                break;
            case '6':
                listAll($recipes);
                echo "Введите ID для удаления: ";
                $id = (int) trim(fgets(STDIN));
                if (deleteRecipe($recipes, $id)) {
                    echo color("✅ Рецепт удалён.\n", '32');
                } else {
                    echo color("❌ Рецепт не найден.\n", '31');
                }
                break;
            case '7':
                listAll($recipes);
                echo "Введите ID для редактирования: ";
                $id = (int) trim(fgets(STDIN));
                echo "Какое поле редактировать (name, category, ingredients, instructions, prep_time, difficulty): ";
                $field = trim(fgets(STDIN));
                echo "Новое значение: ";
                $value = trim(fgets(STDIN));
                if (editRecipe($recipes, $id, $field, $value)) {
                    echo color("✅ Рецепт обновлён.\n", '32');
                } else {
                    echo color("❌ Не удалось обновить.\n", '31');
                }
                break;
            case '8':
                listAll($recipes);
                echo "Введите ID для оценки: ";
                $id = (int) trim(fgets(STDIN));
                echo "Оценка (1-5): ";
                $rating = (int) trim(fgets(STDIN));
                if ($rating >= 1 && $rating <= 5) {
                    if (rateRecipe($recipes, $id, $rating)) {
                        echo color("✅ Оценка добавлена.\n", '32');
                    } else {
                        echo color("❌ Рецепт не найден.\n", '31');
                    }
                } else {
                    echo color("❌ Оценка должна быть от 1 до 5.\n", '31');
                }
                break;
            case '9': stats($recipes); break;
            case '10':
                echo "1. Экспорт в JSON\n";
                echo "2. Экспорт в CSV\n";
                echo "Выберите формат: ";
                $sub = trim(fgets(STDIN));
                if ($sub == '1') {
                    exportJSON($recipes);
                } elseif ($sub == '2') {
                    exportCSV($recipes);
                } else {
                    echo color("Неверный выбор.\n", '31');
                }
                break;
            case '11':
                echo "До свидания!\n";
                exit(0);
            default:
                echo color("Неверный выбор.\n", '31');
        }
    }
}

main();
?>
