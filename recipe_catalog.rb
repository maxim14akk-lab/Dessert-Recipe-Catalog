# recipe_catalog.rb — Ruby версия

require 'json'
require 'csv'

DATA_FILE = 'recipes.json'

class Recipe
  attr_accessor :id, :name, :category, :ingredients, :instructions, :prep_time, :difficulty, :rating, :ratings_count

  def initialize(id, name, category, ingredients, instructions, prep_time, difficulty, rating = 0, ratings_count = 0)
    @id = id
    @name = name
    @category = category
    @ingredients = ingredients
    @instructions = instructions
    @prep_time = prep_time
    @difficulty = difficulty
    @rating = rating
    @ratings_count = ratings_count
  end

  def to_h
    { id: @id, name: @name, category: @category, ingredients: @ingredients,
      instructions: @instructions, prep_time: @prep_time, difficulty: @difficulty,
      rating: @rating, ratings_count: @ratings_count }
  end

  def self.from_h(h)
    new(h[:id], h[:name], h[:category], h[:ingredients], h[:instructions],
        h[:prep_time], h[:difficulty], h[:rating] || 0, h[:ratings_count] || 0)
  end

  def add_rating(rating)
    total = @rating * @ratings_count + rating
    @ratings_count += 1
    @rating = total / @ratings_count.to_f
  end
end

class RecipeCatalog
  attr_reader :recipes

  def initialize
    @recipes = []
    load
  end

  def load
    if File.exist?(DATA_FILE)
      begin
        data = JSON.parse(File.read(DATA_FILE), symbolize_names: true)
        @recipes = data.map { |r| Recipe.from_h(r) }
      rescue
        @recipes = []
      end
    end
  end

  def save
    File.write(DATA_FILE, JSON.pretty_generate(@recipes.map(&:to_h)))
  end

  def add(name, category, ingredients, instructions, prep_time, difficulty)
    id = @recipes.size + 1
    @recipes << Recipe.new(id, name, category, ingredients, instructions, prep_time, difficulty)
    save
    id
  end

  def list_all
    if @recipes.empty?
      puts "\e[33mКаталог пуст.\e[0m"
      return
    end
    printf "\e[36m%-4s %-25s %-12s %-8s %-10s %-8s\e[0m\n", "ID", "Название", "Категория", "Время", "Сложность", "Рейтинг"
    puts "-" * 75
    @recipes.each do |r|
      rating_str = r.ratings_count > 0 ? "#{r.rating.round(1)}★" : "—"
      diff_color = r.difficulty == "лёгкая" ? "\e[32m" : r.difficulty == "средняя" ? "\e[33m" : "\e[31m"
      printf "%-4d %-25s %-12s %-8d %s%-10s\e[0m %-8s\n", r.id, r.name, r.category, r.prep_time, diff_color, r.difficulty, rating_str
    end
  end

  def search(query)
    query = query.downcase
    results = @recipes.select do |r|
      r.name.downcase.include?(query) || r.ingredients.any? { |i| i.downcase.include?(query) }
    end
    if results.empty?
      puts "\e[33mНичего не найдено.\e[0m"
      return
    end
    results.each { |r| puts "#{r.id}: #{r.name} | #{r.category} | #{r.prep_time} мин | #{r.difficulty} | #{r.rating.round(1)}★" }
  end

  def filter_by(field, value)
    results = case field
    when 'category'
      @recipes.select { |r| r.category.casecmp(value).zero? }
    when 'difficulty'
      @recipes.select { |r| r.difficulty.casecmp(value).zero? }
    when 'prep_time'
      @recipes.select { |r| r.prep_time <= value.to_i }
    else
      []
    end
    if results.empty?
      puts "\e[33mНет рецептов, соответствующих фильтру.\e[0m"
      return
    end
    results.each { |r| puts "#{r.id}: #{r.name} | #{r.category} | #{r.prep_time} мин | #{r.difficulty}" }
  end

  def sort_by(field, reverse)
    case field
    when 'name'
      @recipes.sort_by! { |r| r.name }
    when 'prep_time'
      @recipes.sort_by! { |r| r.prep_time }
    when 'rating'
      @recipes.sort_by! { |r| r.rating }
    else
      puts "\e[31mНеверное поле для сортировки.\e[0m"
      return
    end
    @recipes.reverse! if reverse
    list_all
  end

  def delete(id)
    found = @recipes.find { |r| r.id == id }
    if found
      @recipes.delete(found)
      save
      true
    else
      false
    end
  end

  def edit(id, field, value)
    recipe = @recipes.find { |r| r.id == id }
    return false unless recipe
    case field
    when 'name' then recipe.name = value
    when 'category' then recipe.category = value
    when 'ingredients' then recipe.ingredients = value.split(',').map(&:strip)
    when 'instructions' then recipe.instructions = value
    when 'prep_time' then recipe.prep_time = value.to_i
    when 'difficulty' then recipe.difficulty = value
    else return false
    end
    save
    true
  end

  def rate(id, rating)
    recipe = @recipes.find { |r| r.id == id }
    return false unless recipe
    recipe.add_rating(rating)
    save
    true
  end

  def stats
    if @recipes.empty?
      puts "Нет данных."
      return
    end
    categories = Hash.new(0)
    difficulties = Hash.new(0)
    total_rating = 0.0
    rated_count = 0
    @recipes.each do |r|
      categories[r.category] += 1
      difficulties[r.difficulty] += 1
      if r.ratings_count > 0
        total_rating += r.rating
        rated_count += 1
      end
    end
    avg_rating = rated_count > 0 ? total_rating / rated_count : 0
    puts "\e[36m📊 Статистика:\e[0m"
    puts "  Всего рецептов: #{@recipes.size}"
    puts "  Средний рейтинг: #{avg_rating.round(1)}★ (из #{rated_count} оценённых)"
    puts "  По категориям:"
    categories.sort_by { |_, c| -c }.each { |c, count| puts "    #{c}: #{count}" }
    puts "  По сложности:"
    difficulties.each { |d, count| puts "    #{d}: #{count}" }
  end

  def export_json(filename = 'recipes_export.json')
    File.write(filename, JSON.pretty_generate(@recipes.map(&:to_h)))
    puts "\e[32m💾 Экспорт JSON: #{filename}\e[0m"
  end

  def export_csv(filename = 'recipes_export.csv')
    if @recipes.empty?
      puts "\e[33mНет данных для экспорта.\e[0m"
      return
    end
    CSV.open(filename, 'w') do |csv|
      csv << ["ID", "Название", "Категория", "Ингредиенты", "Инструкция", "Время", "Сложность", "Рейтинг"]
      @recipes.each do |r|
        csv << [r.id, r.name, r.category, r.ingredients.join(', '), r.instructions, r.prep_time, r.difficulty, r.rating.round(1)]
      end
    end
    puts "\e[32m💾 Экспорт CSV: #{filename}\e[0m"
  end
end

def main
  catalog = RecipeCatalog.new
  loop do
    puts "\n\e[36m🍰 Dessert Recipe Catalog (Ruby)\e[0m"
    puts "1. Добавить рецепт"
    puts "2. Показать все рецепты"
    puts "3. Поиск рецептов"
    puts "4. Фильтрация"
    puts "5. Сортировка"
    puts "6. Удалить рецепт"
    puts "7. Редактировать рецепт"
    puts "8. Оценить рецепт"
    puts "9. Статистика"
    puts "10. Экспорт"
    puts "11. Выход"
    print "Выберите действие: "
    choice = gets.chomp
    case choice
    when "1"
      print "Название: "
      name = gets.chomp
      print "Категория (торт, пирожное, печенье, мороженое, другое): "
      category = gets.chomp
      print "Ингредиенты (через запятую): "
      ingredients = gets.chomp.split(',').map(&:strip)
      print "Инструкция: "
      instructions = gets.chomp
      print "Время приготовления (мин): "
      prep_time = gets.chomp.to_i
      print "Сложность (лёгкая/средняя/сложная): "
      difficulty = gets.chomp.downcase
      difficulty = "средняя" unless ["лёгкая", "средняя", "сложная"].include?(difficulty)
      id = catalog.add(name, category, ingredients, instructions, prep_time, difficulty)
      puts "\e[32m✅ Рецепт добавлен (ID: #{id})\e[0m"
    when "2"
      catalog.list_all
    when "3"
      print "Введите запрос (название или ингредиент): "
      query = gets.chomp
      catalog.search(query)
    when "4"
      puts "Фильтровать по: category, difficulty, prep_time"
      print "Поле: "
      field = gets.chomp.downcase
      print "Значение: "
      value = gets.chomp
      catalog.filter_by(field, value)
    when "5"
      puts "Сортировать по: name, prep_time, rating"
      print "Поле: "
      field = gets.chomp.downcase
      print "По убыванию? (y/n): "
      reverse = gets.chomp.downcase == 'y'
      catalog.sort_by(field, reverse)
    when "6"
      catalog.list_all
      print "Введите ID для удаления: "
      id = gets.chomp.to_i
      if catalog.delete(id)
        puts "\e[32m✅ Рецепт удалён.\e[0m"
      else
        puts "\e[31m❌ Рецепт не найден.\e[0m"
      end
    when "7"
      catalog.list_all
      print "Введите ID для редактирования: "
      id = gets.chomp.to_i
      print "Какое поле редактировать (name, category, ingredients, instructions, prep_time, difficulty): "
      field = gets.chomp.downcase
      print "Новое значение: "
      value = gets.chomp
      if catalog.edit(id, field, value)
        puts "\e[32m✅ Рецепт обновлён.\e[0m"
      else
        puts "\e[31m❌ Не удалось обновить.\e[0m"
      end
    when "8"
      catalog.list_all
      print "Введите ID для оценки: "
      id = gets.chomp.to_i
      print "Оценка (1-5): "
      rating = gets.chomp.to_i
      if (1..5).include?(rating)
        if catalog.rate(id, rating)
          puts "\e[32m✅ Оценка добавлена.\e[0m"
        else
          puts "\e[31m❌ Рецепт не найден.\e[0m"
        end
      else
        puts "\e[31m❌ Оценка должна быть от 1 до 5.\e[0m"
      end
    when "9"
      catalog.stats
    when "10"
      puts "1. Экспорт в JSON"
      puts "2. Экспорт в CSV"
      print "Выберите формат: "
      sub = gets.chomp
      if sub == "1"
        catalog.export_json
      elsif sub == "2"
        catalog.export_csv
      else
        puts "\e[31mНеверный выбор.\e[0m"
      end
    when "11"
      puts "До свидания!"
      break
    else
      puts "\e[31mНеверный выбор.\e[0m"
    end
  end
end

main if __FILE__ == $0
